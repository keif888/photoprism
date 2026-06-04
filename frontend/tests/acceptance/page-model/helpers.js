import { ClientFunction } from "testcafe";
import testcafeconfig from "../../testcafeconfig.json";

export const showLogs = process.env.SHOW_LOGS === "true";

// getTopElement will return details on what is on top of a selector.
// Useful when the standard output from testcafe warning is insufficient to identify the obstruction.
export const getTopElement = ClientFunction((selectorFn) => {
  const el = selectorFn();
  const rect = el.getBoundingClientRect();
  const centerX = rect.left + rect.width / 2;
  const centerY = rect.top + rect.height / 2;

  // Returns the topmost element at the exact center of the target
  const topEl = document.elementFromPoint(centerX, centerY);

  return {
    tagName: topEl.tagName.toLowerCase(),
    className: topEl.className,
    id: topEl.id,
    innerText: topEl.innerText,
    innerHTML: topEl.innerHTML,
    outerHTML: topEl.outerHTML
  };
});

// clickIfVisible speeds up testcafe's click by shortening the wait time when an item is visible.
// Should only be used for cases where an overlay obscures a selection, but it will not impact the testing.
export async function clickIfVisible(t, sel) {
  if (await sel.visible) {
    await t.click(sel.with({ timeout: 150 }));
  } else {
    await t.click(sel);
  }
}

export function logMessage(message) {
    if (showLogs) {
        const now = new Date();
        console.log(now.toISOString() + " " + message);
    }
}

export function logTime(key) {
  showLogs && console.time(key);
}

export function logTimeEnd(key) {
  showLogs && console.timeEnd(key);
}

// helperRequest performs in a similar way to t.request, but uses fetch instead.
// Currently supports public tests.
// Must be provided with url and method objects
// Example :
// {
//   url: `${testcafeconfig.api}albums`,
//   method: 'get',
//   params: {
//     count: Limit,
//     offset: xOffset
//   }
// }
async function helperRequest(requestOptions) {
  const rURL = new URL(requestOptions.url);
  if (requestOptions.params) {
    Object.entries(requestOptions.params).forEach(([key, value]) => {
      rURL.searchParams.append(key, value);
    });
  }
  
  let rOption = {
    method: requestOptions.method
  }

  if (requestOptions.body) {
    rOption.body = requestOptions.body;
  }

  const response = await fetch(rURL, rOption);

  let headers = {};
  if (response.headers) {
    response.headers.forEach((value, key) => {
      headers[key] = value;
    });
  }

  if (response.ok) {
    const result = await response.json();
    return {
      status: response.status,
      statusText: response.statusText,
      headers: headers,
      body: result
    };
  } else {
    try {
      const result = await response.json();
      return {
        status: response.status,
        statusText: response.statusText,
        headers: headers,
        body: result
      };
    } catch {
      return response;
    }
  }
}

// getAllAlbumPhotosV2 uses Fetch for get all the basic photo information of photos that are in the album identified by it's UID (albumUID).
async function getAllAlbumPhotosV2(albumUID) {
  let photos = [];
  const Limit = 50;
  let xCount = Limit;
  let xOffset = 0;

  while (xCount === Limit) {
    const photoApiResponse = await helperRequest({
      url: `${testcafeconfig.api}photos`,
      params: {
        count: Limit,
        offset: xOffset,
        s: albumUID
      }
    });
    xOffset += Limit;
    if (photoApiResponse.status === 200) {
      xCount = Number(photoApiResponse.headers["x-count"]);
      photos.push(...photoApiResponse.body);
    } else {
      xCount = 0;
      throw new Error(`getAllPhotos failed with status ${photoApiResponse.status} and ${photoApiResponse.statusText}`);
    }
  }
  return photos;
}

// helperBeforeFixture will setup the context for all the tests, by taking a snapshot of the before fixture state.
export async function helperBeforeFixture(ctx) {
  logMessage("helperBeforeFixture");
  let snapshotAlbums = [];
  let snapshotLabels = [];
  let snapshotPhotos = [];
  // ToDo: take a snapshot of all the albums, photos and labels.
  // Although it is possible to take a snapshot of all the albums, photos and labels, 
  // it will take a lot of API calls to achieve this.  
  // Photos: 

  // Labels is just the base query as there are not archived labels or details that can be queried (10ms).
  // get /api/v1/labels?count=150

  let searchApiResponse;
  const Limit = 50;
  let xCount = Limit;
  let xOffset = 0;
  // Snapshot all albums, and the photos in them
  while (xCount === Limit) {
    searchApiResponse = await helperRequest({
      url: `${testcafeconfig.api}albums`,
      method: 'get',
      params: {
        count: Limit,
        offset: xOffset
      }
    });
    xOffset += Limit;
    if (searchApiResponse.status === 200) {
      xCount = Number(searchApiResponse.headers["x-count"]);
      for (const album of searchApiResponse.body) {
        const photos = await getAllAlbumPhotosV2(album.UID);
        const rAlbum = {
          "uid": album.UID,
          "data": album,
          "photos": photos
        }
        snapshotAlbums.push(rAlbum);
      }
    } else {
      const msg = "helperBeforeFixture gather albums " + JSON.stringify(searchApiResponse);
      logMessage(msg);
      throw new Error(msg);
    }
  }

  xCount = Limit;
  xOffset = 0;
  // Snapshot all labels (the photo snapshot captures which labels are for which photos)
  while (xCount === Limit) {
    searchApiResponse = await helperRequest({
      url: `${testcafeconfig.api}labels`,
      method: 'get',
      params: {
        count: Limit,
        offset: xOffset,
        all: true
      }
    });
    xOffset += Limit;
    if (searchApiResponse.status === 200) {
      xCount = searchApiResponse.body.length;  // labels does not return an x-count as at 2026-06-04.
      for (const label of searchApiResponse.body) {
        const rLabel = {
          "uid": label.UID,
          "data": label
        }
        snapshotLabels.push(rLabel);
      }
    } else {
      const msg = "helperBeforeFixture gather labels " + JSON.stringify(searchApiResponse);
      logMessage(msg);
      throw new Error(msg);
    }
  }

  // There are 110 photos with files in Acceptance database
  // primary:true public:false - 104 (includes review photos)
  // primary:true archived:true - 6 (archived:true overrides public:false, so only archived photos are returned.)
  // Snapshot non archived photos, getting all their details (2nd call per photo required)
  xCount = Limit;
  xOffset = 0;
  while (xCount === Limit) {
    searchApiResponse = await helperRequest({
      url: `${testcafeconfig.api}photos`,
      method: 'get',
      params: {
        count: Limit,
        offset: xOffset,
        q: "primary:true public:false"
      }
    });
    xOffset += Limit;
    if (searchApiResponse.status === 200) {
      xCount = Number(searchApiResponse.headers["x-count"]);
      for (const photo of searchApiResponse.body) {
        const photoApiResponse = await helperRequest({
          url: `${testcafeconfig.api}photos/${photo.UID}`,
          method: 'get'
        });
        if (photoApiResponse.status === 200) {
          const rPhoto = {
            "uid": photo.UID,
            "data": photoApiResponse.body
          }          
          snapshotPhotos.push(rPhoto);
        } else {
          const msg = "helperBeforeFixture query photo " + JSON.stringify(photoApiResponse);
          logMessage(msg);
          throw new Error(msg);
        }
      }
    } else {
      const msg = "helperBeforeFixture gather photos " + JSON.stringify(searchApiResponse);
      logMessage(msg);
      throw new Error(msg);
    }
  }
  // Snapshot archived photos, getting all their details (2nd call per photo required)
  xCount = Limit;
  xOffset = 0;
  while (xCount === Limit) {
    searchApiResponse = await helperRequest({
      url: `${testcafeconfig.api}photos`,
      method: 'get',
      params: {
        count: Limit,
        offset: xOffset,
        q: "primary:true archived:true"
      }
    });
    xOffset += Limit;
    if (searchApiResponse.status === 200) {
      xCount = Number(searchApiResponse.headers["x-count"]);
      for (const photo of searchApiResponse.body) {
        const photoApiResponse = await helperRequest({
          url: `${testcafeconfig.api}photos/${photo.UID}`,
          method: 'get'
        });
        if (photoApiResponse.status === 200) {
          const rPhoto = {
            "uid": photo.UID,
            "data": photoApiResponse.body
          }          
          snapshotPhotos.push(rPhoto);
        } else {
          const msg = "helperBeforeFixture query photo " + JSON.stringify(photoApiResponse);
          logMessage(msg);
          throw new Error(msg);
        }
      }
    } else {
      const msg = "helperBeforeFixture gather photos " + JSON.stringify(searchApiResponse);
      logMessage(msg);
      throw new Error(msg);
    }
}

  // Store all the snapshots for use in the beforeEach/afterEach functions
  ctx.snapshots = {
    "snapshotAlbums": snapshotAlbums,
    "snapshotLabels": snapshotLabels,
    "snapshotPhotos": snapshotPhotos
  }
}

// helperBeforeEach will setup the context for test reversion
export async function helperBeforeEachV1(t) {
  logMessage("helperBeforeEach");
  t.ctx.testChanges = {
    "revertAlbums": [],
    "revertPhotos": [],
    "removeAlbums": [],
    "removeLabels": [],
    "removeLabelFromPhotos": []
  }
}

async function getAllAlbumPhotos(t, albumUID) {
  let photos = [];
  const Limit = 50;
  let xCount = Limit;
  let xOffset = 0;

  while (xCount === Limit) {
    const photoApiResponse = await t.request({
      url: `${testcafeconfig.api}photos`,
      params: {
        count: Limit,
        offset: xOffset,
        s: albumUID
      }
    });
    xOffset += Limit;
    if (photoApiResponse.status === 200) {
      xCount = Number(photoApiResponse.headers["x-count"]);
      photos.push(...photoApiResponse.body);
    } else {
      xCount = 0;
      throw new Error(`getAllPhotos failed with status ${photoApiResponse.status} and ${photoApiResponse.statusText}`);
    }
  }
  return photos;
}

export async function helperBeforeEach(t) {
  logMessage("helperBeforeEach");

    let startTimestamp = new Date();
  startTimestamp.setMilliseconds(0);
  let helperFailures = [];
  // let snapshotAlbums = [];
  // ToDo: take a snapshot of all the albums, photos and labels.
  // Although it is possible to take a snapshot of all the albums, photos and labels, 
  // it will take a lot of API calls to achieve this.  
  // Photos: 
  // Page through get /api/v1/photos without a query (40ms), and call get /api/v1/photos/{uid} (10ms) for each photo found, and store the result
  // Page through get /api/v1/photos with query archive:true (40ms), and call get /api/v1/photos/{uid} for each photo found, and store the result
  // There are 114 photos, so this will take at least 116 api calls, assuming that the page size is 150.  Execution time ~1.5s.
  // Albums is just the base query as there are not archived albums or details that can be queried (10ms).
  // get /api/v1/albums?count=150
  // Labels is just the base query as there are not archived labels or details that can be queried (10ms).
  // get /api/v1/labels?count=150

  // let searchApiResponse;
  // const Limit = 50;
  // let xCount = 50;
  // let xOffset = 0;
  // // Find Albums to Revert
  // while (xCount === Limit) {
  //   searchApiResponse = await t.request({
  //     url: `${testcafeconfig.api}albums`,
  //     method: 'get',
  //     params: {
  //       count: Limit,
  //       offset: xOffset
  //     }
  //   });
  //   xOffset += Limit;
  //   if (searchApiResponse.status === 200) {
  //     xCount = Number(searchApiResponse.headers["x-count"]);
  //     for (const album of searchApiResponse.body) {
  //       const photos = await getAllAlbumPhotos(t, album.UID);
  //       const rAlbum = {
  //         "uid": album.UID,
  //         "data": album,
  //         "photos": photos
  //       }
  //       snapshotAlbums.push(rAlbum);
  //     }
  //   } else {
  //     const msg = "helperBeforeEach gather albums " + JSON.stringify(apiResponse);
  //     logMessage(msg);
  //     helperFailures.push(msg);
  //   }
  // }

  t.ctx.testChanges = {
    "startTimestamp": startTimestamp.toISOString(),
    "revertAlbums": [],
    "revertPhotos": [],
    "removeAlbums": [],
    "removeLabels": [],
    "removeLabelFromPhotos": []
  }
  await t.expect(helperFailures).eql([]);
}


// this function stores the current information about a photo that will need to be reverted.
export async function helperRevertAlbum (t, uid) {
  logMessage(`helperRevertAlbum (t, ${uid})`);
  await saveRevertAlbum(t, uid);
}

async function saveRevertAlbum (t, uid) {
  // Get the album details
  const apiResponse = await t.request(`${testcafeconfig.api}albums/${uid}`);
  const photoApiResponse = await getAllAlbumPhotos(t, uid);
  const revertAlbum = {
    "uid": uid,
    "data": apiResponse.body,
    "photos": photoApiResponse.body
  }
  if (!t.ctx.testChanges.revertAlbums.some(ra => ra.uid === uid)) {
    t.ctx.testChanges.revertAlbums.push(revertAlbum);
  }
}

// this function stores the current information about a photo that will need to be reverted.
// Known issues after executing helperAfterEach:
// Automatically generated titles may be updated to match new format (over old in acceptance data), or reflect labels reverted
// Labels will change type and uncertainty if they were not manual and need to be reverted
// Updated timestamps will change
// Can not restack a file that has been unstacked from a photo
// Can not undelete a file that has been deleted from a photo
// Can NOT revert a photo that has Quality < 3, because the quality will be 3 after reversion due to edits applied.
export async function helperRevertPhoto (t, uid) {
  logMessage(`helperRevertPhoto (t, ${uid})`);
  const apiResponse = await t.request(`${testcafeconfig.api}photos/${uid}`);
  const revertPhoto = {
    "uid": uid,
    "data": apiResponse.body
  }
  if (apiResponse.body.Quality < 3) {
    logMessage(`Photo was not added to list of photos to revert as it is in review status.`);
  } else {
    if (!t.ctx.testChanges.revertPhotos.some(rp => rp.uid === uid)) {
      t.ctx.testChanges.revertPhotos.push(revertPhoto);
    }
  }
}

// this function stores the need to remove an album.
export async function helperRemoveAlbum (t, how, id) {
  logMessage(`helperRemoveAlbum (t, ${how}, ${id})`);
  if (how === "name") {
    const removeAlbum = {
      "name": id,
      "uid": "name"
    }
    if (!t.ctx.testChanges.removeAlbums.some(ra => ra.name === id)) {
      t.ctx.testChanges.removeAlbums.push(removeAlbum);
    }
  } else {
    const removeAlbum = {
      "name": "uid",
      "uid": id
    }
    if (!t.ctx.testChanges.removeAlbums.some(ra => ra.uid === id)) {
      t.ctx.testChanges.removeAlbums.push(removeAlbum);
    }
  }
}

// this function stores the need to remove a label from a photo.
export async function helperRemoveLabelFromPhotos (t, labelUid, photoUid) {
  logMessage(`helperRemoveLabelFromPhotos (t, ${labelUid}, ${photoUid})`);
  if (photoUid.length > 0) {
    if (!isNaN(+labelUid)) {
      const removeLabelFromPhoto = {
        "labelUid": labelUid,
        "photoUid" : photoUid
      }
      if (!t.ctx.testChanges.removeLabelFromPhotos.some(rlfp => rlfp.labelUid === labelUid && rlfp.photoUid === photoUid)) {
        t.ctx.testChanges.removeLabelFromPhotos.push(removeLabelFromPhoto);
      }
    } else {
      const apiResponse = await t.request(`${testcafeconfig.api}labels?count=1&q=uid:${labelUid}`);
      if (apiResponse.status === 200) {
        const removeLabelFromPhoto = {
          "labelUid": apiResponse.body[0].ID,
          "photoUid" : photoUid
        }
        if (!t.ctx.testChanges.removeLabelFromPhotos.some(rlfp => rlfp.labelUid === labelUid && rlfp.photoUid === photoUid)) {
          t.ctx.testChanges.removeLabelFromPhotos.push(removeLabelFromPhoto);
        }
      }
    }
  }
}

// this function stores the need to remove a label from ALL photos.
// Please note that this may leave invalid titles on photos.
export async function helperRemoveLabel (t, name) {
  logMessage(`helperRemoveLabel (t, ${name})`);
  const removeLabel = {
    "name": name
  }
  if (!t.ctx.testChanges.removeLabels.some(rl => rl.name === name)) {
    t.ctx.testChanges.removeLabels.push(removeLabel);
  }
}

// helperDetermineChangedItems detects album, label and/or photo changes, and populates the list of changed items for reversion/removal.
async function helperDetermineChangedItems(t) {
  const beforeTimestamp = new Date(t.ctx.testChanges.startTimestamp);
  // Start remove these after the other helpers are removed.
  t.ctx.testChanges.revertAlbums.length = 0;
  t.ctx.testChanges.removeAlbums.length = 0;
  t.ctx.testChanges.removeLabels.length = 0;
  t.ctx.testChanges.revertPhotos.length = 0;
  t.ctx.testChanges.removeLabelFromPhotos.length = 0;
  // End remove these after the other helpers are removed.
  let searchApiResponse;
  const Limit = 50;
  // Find Albums to Revert/Remove due to changes or being added.
  let xCount = Limit;
  let xOffset = 0;
  let foundAlbums = [];
  while (xCount === Limit) {
    searchApiResponse = await t.request({
      url: `${testcafeconfig.api}albums`,
      method: 'get',
      params: {
        count: Limit,
        offset: xOffset
      }
    });
    xOffset += Limit;
    if (searchApiResponse.status === 200) {
      xCount = Number(searchApiResponse.headers["x-count"]);
      foundAlbums.push(...searchApiResponse.body.map(album => album.UID));
      const albums = searchApiResponse.body.filter((album) => {
        return new Date(album.CreatedAt) >= beforeTimestamp || new Date(album.UpdatedAt) >= beforeTimestamp
      });
      for (const album of albums) {
        if (new Date(album.CreatedAt) >= beforeTimestamp) {
          const removeAlbum = {
            "name": "uid",
            "uid": album.UID
          }
          if (!t.ctx.testChanges.removeAlbums.some(ra => ra.uid === album.UID)) {
            t.ctx.testChanges.removeAlbums.push(removeAlbum);
          }
        } else {
          // await saveRevertAlbum(t, album.UID);
          const rAlbum = t.fixtureCtx.snapshots.snapshotAlbums.find((a) => a.uid === album.UID);
          if (rAlbum) {
            t.ctx.testChanges.revertAlbums.push(rAlbum);
          } else {
            const msg = `helperDetermineChangedItems revert albums (1) ${album.UID} is missing from snapshot albums`;
            logMessage(msg);
            return msg;
          }
        }
      }
    } else {
      const msg = "helperDetermineChangedItems revert albums (2) " + JSON.stringify(apiResponse);
      logMessage(msg);
      return msg;
    }
  }
  // Find albums that have been deleted, and flag them for reversion.
  for (const album of t.fixtureCtx.snapshots.snapshotAlbums) {
    if (!foundAlbums.includes(album.uid)) {
      t.ctx.testChanges.revertAlbums.push(album);
    }
  }

  // Find Labels to remove
  xCount = Limit;
  xOffset = 0;
  while (xCount === Limit) {
    searchApiResponse = await helperRequest({
      url: `${testcafeconfig.api}labels`,
      method: 'get',
      params: {
        count: Limit,
        offset: xOffset,
        all: true
      }
    });
    xOffset += Limit;
    if (searchApiResponse.status === 200) {
      xCount = searchApiResponse.body.length;
      const labels = searchApiResponse.body.filter((label) => {
        return new Date(label.CreatedAt) >= beforeTimestamp
      });
      for (const label of labels) {
        const removeLabel = {
          "uid": label.UID
        }
        if (!t.ctx.testChanges.removeLabels.some(ra => ra.uid === label.UID)) {
          t.ctx.testChanges.removeLabels.push(removeLabel);
        }
      }
    } else {
      const msg = "helperDetermineChangedItems remove labels " + JSON.stringify(apiResponse);
      logMessage(msg);
      return msg;
    }
  }


  // There are 110 photos with files in Acceptance database
  // primary:true public:false - 104 (includes review photos)
  // primary:true archived:true - 6 (archived:true overrides public:false, so only archived photos are returned.)
  // Find non archived photos with UpdatedAt >= beforeTimestamp
  let foundPhotos = [];
  xCount = Limit;
  xOffset = 0;
  while (xCount === Limit) {
    searchApiResponse = await helperRequest({
      url: `${testcafeconfig.api}photos`,
      method: 'get',
      params: {
        count: Limit,
        offset: xOffset,
        q: "primary:true public:false"
      }
    });
    xOffset += Limit;
    if (searchApiResponse.status === 200) {
      xCount = Number(searchApiResponse.headers["x-count"]);
      foundPhotos.push(...searchApiResponse.body.map(photo => photo.UID));
      const photos = searchApiResponse.body.filter((photo) => {
        return new Date(photo.UpdatedAt) >= beforeTimestamp; 
      });

      for (const photo of photos) {
        const rPhoto = t.fixtureCtx.snapshots.snapshotPhotos.find((a) => a.uid === photo.UID);
        if (rPhoto) {
          t.ctx.testChanges.revertPhotos.push(rPhoto);
        } else {
          const msg = `helperDetermineChangedItems revert photos (1) ${photo.UID} is missing from snapshot photos`;
          logMessage(msg);
          return msg;
        }
      }
    } else {
      const msg = "helperDetermineChangedItems gather photos (1) " + JSON.stringify(searchApiResponse);
      logMessage(msg);
      return msg;
    }
  }
  // Find archived photos with DeletedAt or UpdatedAt >= beforeTimestamp
  xCount = Limit;
  xOffset = 0;
  while (xCount === Limit) {
    searchApiResponse = await helperRequest({
      url: `${testcafeconfig.api}photos`,
      method: 'get',
      params: {
        count: Limit,
        offset: xOffset,
        q: "primary:true archived:true"
      }
    });
    xOffset += Limit;
    if (searchApiResponse.status === 200) {
      xCount = Number(searchApiResponse.headers["x-count"]);
      foundPhotos.push(...searchApiResponse.body.map(photo => photo.UID));
      const photos = searchApiResponse.body.filter((photo) => {
        return new Date(photo.UpdatedAt) >= beforeTimestamp || new Date(photo.DeletedAt) >= beforeTimestamp; 
      });

      for (const photo of photos) {
        const rPhoto = t.fixtureCtx.snapshots.snapshotPhotos.find((a) => a.uid === photo.UID);
        if (rPhoto) {
          t.ctx.testChanges.revertPhotos.push(rPhoto);
        } else {
          const msg = `helperDetermineChangedItems revert photos (2) ${photo.UID} is missing from snapshot photos`;
          logMessage(msg);
          return msg;
        }
      }
    } else {
      const msg = "helperDetermineChangedItems gather photos (2) " + JSON.stringify(searchApiResponse);
      logMessage(msg);
      return msg;
    }
  }

  logMessage(JSON.stringify(t.ctx.testChanges));
  console.log(JSON.stringify(t.ctx.testChanges));
  return '';
}


// This function will undo what the test has done (to the best of it's ability)
// as requested by the helperRemove and helperRevert functions.
export async function helperAfterEach(t) {
  logMessage("helperAfterEach");
  const result = await helperDetermineChangedItems(t);

  let helperFailures = [];
  // logMessage("helperAfterEach Queued Requests " + JSON.stringify(t.ctx.testChanges));
  // Revert Albums state
  // This MAY result in a different UID if the album has been deleted, and it wasn't created by the current user.
  try {
    for (let revertAlbum of t.ctx.testChanges.revertAlbums) {
      let apiResponse = await t.request({
        url: `${testcafeconfig.api}albums/${revertAlbum.uid}`,
        method: 'put',
        body: revertAlbum.data
      });
      if (apiResponse.status === 404) {
        let apiPostResponse = await t.request({
          url: `${testcafeconfig.api}albums`,
          method: 'post',
          body: revertAlbum.data
        });

        if (apiPostResponse.status === 201) { // The Album has been created with a different UID!
          revertAlbum.data.UID = apiPostResponse.body.UID;
          revertAlbum.data.ID = apiPostResponse.body.ID;
          revertAlbum.uid = apiPostResponse.body.UID;
          if (revertAlbum.data.Thumb) {
            revertAlbum.data.ThumbSrc = "manual"; // To revert a thumb it must be manual
          }
          // Updating the NEW album
          apiResponse = await t.request({
            url: `${testcafeconfig.api}albums/${revertAlbum.uid}`,
            method: 'put',
            body: revertAlbum.data
          });
          if (apiResponse.status !== 200 || apiResponse.status === null) { // Ignore Ok
            const msg = "helperAfterEach revert albums (1) " + JSON.stringify(apiResponse);
            logMessage(msg);
            helperFailures.push(msg);
          }
        }
      } else {
        if (apiResponse.status !== 200 || apiResponse.status === null) { // Ignore Ok
          const msg = "helperAfterEach revert albums (2) " + JSON.stringify(apiResponse);
          logMessage(msg);
          helperFailures.push(msg);
        }
      }
      // Restore the photos connections
      const albumPhotoApiResponse = await t.request(`${testcafeconfig.api}photos?count=50&offset=0&s=${revertAlbum.uid}`);

      let photos = [];
      for (const photo of revertAlbum.photos) {
        if (!albumPhotoApiResponse.body.find(ap => ap.UID === photo.UID))
        {
          photos.push(photo.UID);
          logMessage(`Reverting album add photo ${photo.UID}`);
        }
      }
      if (photos.length > 0){
        const photoApiResponse = await t.request({
          url: `${testcafeconfig.api}albums/${revertAlbum.uid}/photos`,
          method: 'post',
          body: { "photos": photos }
        });
        if (photoApiResponse.status !== 200 || photoApiResponse.status === null) { // Ignore Ok
          const msg = "helperAfterEach revert albums photos " + JSON.stringify(photoApiResponse);
          logMessage(msg);
          helperFailures.push(msg);
        }
        // Try updating the album again in case the thumb was from a removed photo.
        if (revertAlbum.data.Thumb) {
          revertAlbum.data.ThumbSrc = "manual"; // To revert a thumb it must be manual
        }

        let apiResponse = await t.request({
          url: `${testcafeconfig.api}albums/${revertAlbum.uid}`,
          method: 'put',
          body: revertAlbum.data
        });
        if (apiResponse.status !== 200 || apiResponse.status === null) { // Ignore Ok
          const msg = "helperAfterEach revert albums with thumb manual " + JSON.stringify(apiResponse);
          logMessage(msg);
          helperFailures.push(msg);
        }
      }
    }
  } catch (e) {
    const errorText = e.errmsg || e.message || "An unknown error occurred";
    helperFailures.push(`revertAlbum threw ${errorText}`);
  }
  
  // Revert Photos state
  // this can not fully restore labels to a photo.
  // if the label has been fully removed, and it matches a keyword, then it will 
  // be restored as a keyword based label.  Otherwise it will be a manual
  // style label.
  try {
    for (const revertPhoto of t.ctx.testChanges.revertPhotos) {
      console.log(revertPhoto.uid);
      // Get current photo status
      let apiResponse = await t.request({
        url: `${testcafeconfig.api}photos/${revertPhoto.uid}`,
        method: 'get'
      });
      if (apiResponse.status !== 200 || apiResponse.status === null) { // Ignore Ok
        const msg = "helperAfterEach revert photo (1) " + JSON.stringify(apiResponse);
        logMessage(msg);
        helperFailures.push(msg);
      }

      if (!revertPhoto.data.DeletedAt && apiResponse.body.DeletedAt) {
        // Need to restore the photo
        const restoreResponse = await t.request({
          url: `${testcafeconfig.api}batch/photos/restore`,
          method: 'post',
          body: {
            "photos": [ revertPhoto.uid ]
          }
        });
        if (restoreResponse.status !== 200 || restoreResponse.status === null) { // Ignore Ok
          const msg = "helperAfterEach revert restore photo " + JSON.stringify(restoreResponse);
          logMessage(msg);
          helperFailures.push(msg);
        }
      }

      // Revert the photo
      apiResponse = await t.request({
        url: `${testcafeconfig.api}photos/${revertPhoto.uid}`,
        method: 'put',
        body: revertPhoto.data
      });
      if (apiResponse.status !== 200 || apiResponse.status === null) { // Ignore Ok
        const msg = "helperAfterEach revert photo (2) " + JSON.stringify(apiResponse);
        logMessage(msg);
        helperFailures.push(msg);
      }

      // Loop through the labels in revertPhoto.data and apiResponse.body to add/remove as needed.
      // Remove
      for (const label of apiResponse.body.Labels) {
        const exists = revertPhoto.data.Labels.some(slug => slug.Label.Slug === label.Label.Slug);
        if (!exists) {
          const labelApiResponse = await t.request({
            url: `${testcafeconfig.api}photos/${revertPhoto.uid}/label/${label.LabelID}`,
            method: 'delete'
          });
          if ((labelApiResponse.status !== 200 && labelApiResponse.status !== 404) || labelApiResponse.status === null ) { // Ignore Ok and not found
            const msg = "helperAfterEach remove label from photo " + JSON.stringify(labelApiResponse);
            logMessage(msg);
            helperFailures.push(msg);
          }
        }
      }
      // Add
      for (const label of revertPhoto.data.Labels) {
        const exists = apiResponse.body.Labels.some(slug => slug.Label.Slug === label.Label.Slug);
        if (!exists) {
          const labelApiResponse = await t.request({
            url: `${testcafeconfig.api}photos/${revertPhoto.uid}/label`,
            method: 'post',
            body: {
                "Description": label.Label.Description,
                "Favorite": label.Label.Favorite,
                "Name": label.Label.Name,
                "Notes": label.Label.Notes,
                "Priority": label.Priority,
                "Thumb": label.Label.Thumb,
                "ThumbSrc": label.ThumbSrc,
                "Uncertainty": label.Label.Uncertainty
            }
          });
          if (labelApiResponse.status !== 200 || labelApiResponse.status === null) { // Ignore Ok
            const msg = "helperAfterEach add label " + JSON.stringify(labelApiResponse);
            logMessage(msg);
            helperFailures.push(msg);
          }
        } else {
          const labelApiResponse = await t.request({
            url: `${testcafeconfig.api}photos/${revertPhoto.uid}/label/${label.LabelID}`,
            method: 'put',
            body: {
                "Uncertainty": 0 // Although this doesn't match the previous number, it forces a manual label back into place.  All that can be done.
            }
          });
          if (labelApiResponse.status !== 200 || labelApiResponse.status === null) { // Ignore Ok
            const msg = "helperAfterEach reset label " + JSON.stringify(labelApiResponse);
            logMessage(msg);
            helperFailures.push(msg);
          }
        }
      }

      // Loop through the Albums in revertPhoto.data and apiResponse.body to add/remove as needed.
      // Remove
      for (const album of apiResponse.body.Albums) {
        const exists = revertPhoto.data.Albums.some(slug => slug.Slug === album.Slug);
        if (!exists) {
          const albumApiResponse = await t.request({
            url: `${testcafeconfig.api}albums/${album.UID}/photos`,
            method: 'delete',
            body: {
              "photos": [ revertPhoto.uid ]
            }
          });
          if (albumApiResponse.status !== 200 || albumApiResponse.status === null) { // Ignore Ok
            const msg = "helperAfterEach delete from album " + JSON.stringify(albumApiResponse);
            logMessage(msg);
            helperFailures.push(msg);
          }
        }
      }
      // Add
      for (const album of revertPhoto.data.Albums) {
        const exists = apiResponse.body.Albums.some(slug => slug.Slug === album.Slug);
        if (!exists) {
          const albumApiResponse = await t.request({
            url: `${testcafeconfig.api}albums/${album.UID}/photos`,
            method: 'post',
            body: {
              "photos": [ revertPhoto.uid ]
            }
          });
          if (albumApiResponse.status !== 200 || albumApiResponse.status === null) { // Ignore Ok
            const msg = "helperAfterEach add to album " + JSON.stringify(albumApiResponse);
            logMessage(msg);
            helperFailures.push(msg);
          }
        }
      }

      // Loop through the files and markers to update as required
      // Invalidate any that shouldn't be there.
      for (const file of apiResponse.body.Files) {
        const rFile = revertPhoto.data.Files.find(fileI => fileI.UID === file.UID)
        if (rFile) {
          for (const marker of file.Markers) {
            const rMarker = rFile.Markers.find(m => m.UID === marker.UID && m.FileUID === marker.FileUID);
            let markerApiResponse;
            if (rMarker) {
              // reset
              markerApiResponse = await t.request({
                url: `${testcafeconfig.api}markers/${rMarker.UID}`,
                method: 'put',
                body: rMarker
              });
            } else {
              // inactivate
              markerApiResponse = await t.request({
                url: `${testcafeconfig.api}markers/${marker.UID}`,
                method: 'put',
                body: {
                  "Invalid":true
                }
              });
            }
            if (markerApiResponse.status !== 200 || markerApiResponse.status === null) { // Ignore Ok
              const msg = "helperAfterEach sync markers (1) file " + marker.FileUID + " marker " + marker.UID + " " + JSON.stringify(markerApiResponse);
              logMessage(msg);
              helperFailures.push(msg);
            }
          }
        } else {
          const msg = `Choosing not to remove file ${file.UID} which has been added, as that will break future tests as the file is physically deleted.`;
          logMessage(msg);
          helperFailures.push(msg);
        }
      }
      for (const file of revertPhoto.data.Files) {
        const cFile = apiResponse.body.Files.find(fileI => fileI.UID === file.UID)
        if (cFile) {
          for (const marker of file.Markers) {
            // Restore the marker whether it is there or not.
            const markerApiResponse = await t.request({
                url: `${testcafeconfig.api}markers/${marker.UID}`,
                method: 'put',
                body: marker
              });
            if (markerApiResponse.status !== 200 || markerApiResponse.status === null) { // Ignore Ok
              const msg = "helperAfterEach sync markers (2)" + JSON.stringify(markerApiResponse);
              logMessage(msg);
              helperFailures.push(msg);
            }
          }
        } else {
          const msg = `Unable to restore file ${file.UID} which has been removed.  Flagging error as this situation indicated that a fila has been deleted from the file system.`;
          logMessage(msg);
          helperFailures.push(msg);
        }
      }

      // Revert any changes to Primary file.
      const originalPrimary = revertPhoto.data.Files.find((element) => element.Primary === true)
      const currentPrimary = apiResponse.body.Files.find((element) => element.Primary === true)
      console.log(originalPrimary);
      console.log(currentPrimary);
      if (originalPrimary && currentPrimary) {
        const originalUID = originalPrimary.UID
        const currentUID = currentPrimary.UID
        if (originalUID !== currentUID) {
          const primaryApiResponse = await t.request({
            url: `${testcafeconfig.api}photos/${revertPhoto.uid}/files/${originalUID}/primary`,
            method: 'post'
          });
          if (primaryApiResponse.status !== 200 || primaryApiResponse.status === null) { // Ignore Ok
            const msg = "helperAfterEach revert photo primary " + JSON.stringify(primaryApiResponse);
            logMessage(msg);
            helperFailures.push(msg);
          }
        }
      }

      // Do the photo again to try the Title again.
      apiResponse = await t.request({
        url: `${testcafeconfig.api}photos/${revertPhoto.uid}`,
        method: 'put',
        body: revertPhoto.data
      });
      if (apiResponse.status !== 200 || apiResponse.status === null) { // Ignore Ok
        const msg = "helperAfterEach revert photo again " + JSON.stringify(apiResponse);
        logMessage(msg);
        helperFailures.push(msg);
      }
      if (revertPhoto.data.DeletedAt && !apiResponse.body.DeletedAt) {
        // Need to archive the photo
        const archiveResponse = await t.request({
          url: `${testcafeconfig.api}batch/photos/archive`,
          method: 'post',
          body: {
            "photos": [ revertPhoto.uid ]
          }
        });
        if (archiveResponse.status !== 200 || archiveResponse.status === null) { // Ignore Ok
          const msg = "helperAfterEach revert archive photo " + JSON.stringify(archiveResponse);
          logMessage(msg);
          helperFailures.push(msg);
        }
      }

    }
  } catch (e) {
    const errorText = e.errmsg || e.message || "An unknown error occurred";
    helperFailures.push(`revertAlbum threw ${errorText}`);
  }

  // Remove albums
  try {
    for (const removeAlbum of t.ctx.testChanges.removeAlbums) {
      let listApiResponse;
      if (removeAlbum.uid === "name") {
        listApiResponse = await t.request({
          url: `${testcafeconfig.api}albums`,
          method: 'get',
          params: {
            count: 10,
            q: `${removeAlbum.name} type:album`
          }
        });
      } else {
        listApiResponse = await t.request({
          url: `${testcafeconfig.api}albums`,
          method: 'get',
          params: {
            count: 10,
            q: `uid:${removeAlbum.uid}`
          }
        });
      }
      if (listApiResponse.status !== 200 || listApiResponse.status === null) { // Ignore Ok
        const msg = "helperAfterEach list albums " + JSON.stringify(listApiResponse);
        logMessage(msg);
        helperFailures.push(msg);
      }
      for (const album of listApiResponse.body) {
        const apiResponse = await t.request({
          url: `${testcafeconfig.api}albums/${album.UID}`,
          method: 'delete',
          params: {
              force: true
          }
        });
        if (apiResponse.status !== 200 || apiResponse.status === null && apiResponse.status !== 404) { // Ignore Ok and not found
          const msg = "helperAfterEach delete album " + JSON.stringify(apiResponse);
          logMessage(msg);
          helperFailures.push(msg);
        }
      }
    }
  } catch (e) {
    const errorText = e.errmsg || e.message || "An unknown error occurred";
    helperFailures.push(`removeAlbums threw ${errorText}`);
  }

  // // Remove Labels from Photos
  // try {
  //   for (const removeLabelFromPhoto of t.ctx.testChanges.removeLabelFromPhotos) {
  //     const apiResponse = await t.request({
  //       url: `${testcafeconfig.api}photos/${removeLabelFromPhoto.photoUid}/label/${removeLabelFromPhoto.labelUid}`,
  //       method: 'delete'
  //     });
  //     if ((apiResponse.status !== 200 && apiResponse.status !== 404) || apiResponse.status === null ) { // Ignore Ok and not found
  //       const msg = "helperAfterEach remove label from photo " + JSON.stringify(archiveResponse);
  //       logMessage(msg);
  //       helperFailures.push(msg);
  //     }
  //   }
  // } catch (e) {
  //   const errorText = e.errmsg || e.message || "An unknown error occurred";
  //   helperFailures.push(`removeLabelsFromPhotos threw ${errorText}`);
  // }

  // Remove Labels
  try {
    if (t.ctx.testChanges.removeLabels.length > 0) {
      let labels = [];
      for (const removeLabel of t.ctx.testChanges.removeLabels) {
        // const listApiResponse = await t.request({
        //   url: `${testcafeconfig.api}labels`,
        //   method: 'get',
        //   params: {
        //     count: 10,
        //     q: `${removeLabel.name}`
        //   }
        // });
        // if (listApiResponse.status !== 200 || listApiResponse.status === null) { // Ignore Ok
        //   const msg = "helperAfterEach get labels " + JSON.stringify(listApiResponse);
        //   logMessage(msg);
        //   helperFailures.push(msg);
        // }
        // for (const label of listApiResponse.body) {
        //   labels.push(label.UID);
        // }
        labels.push(removeLabel.uid);
      }
      if (labels.length > 0) {
        const apiResponse = await t.request({
          url: `${testcafeconfig.api}batch/labels/delete`,
          method: 'post',
          body: {
            "labels": labels
          }
        });
        if (apiResponse.status !== 200 || apiResponse.status === null) { // Ignore Ok
          const msg = "helperAfterEach delete labels " + JSON.stringify(apiResponse);
          logMessage(msg);
          helperFailures.push(msg);
        }
      }
    }
  } catch (e) {
    const errorText = e.errmsg || e.message || "An unknown error occurred";
    helperFailures.push(`removeLabels threw ${errorText}`);
  }

  // Error if there were any API or try/catch failures.
  await t.expect(helperFailures).eql([]);
}