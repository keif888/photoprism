import { ClientFunction } from "testcafe";
import testcafeconfig from "../../testcafeconfig.json";

export const showLogs = process.env.SHOW_LOGS == "true";

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
        var now = new Date();
        console.log(now.toISOString() + " " + message);
    }
}

export function logTime(key) {
  showLogs && console.time(key);
}

export function logTimeEnd(key) {
  showLogs && console.timeEnd(key);
}

// helperBeforeEach will setup the context for test reversion
export function helperBeforeEach(t) {
  t.ctx.testChanges = {
    "revertAlbums": [],
    "revertPhotos": [],
    "removeAlbums": [],
    "removeLabels": [],
    "removeLabelFromPhotos": []
  }
}

// this function stores the current information about a photo that will need to be reverted.
export async function helperRevertAlbum (t, uid) {
  const apiResponse = await t.request(`${testcafeconfig.api}albums/${uid}`);
  const revertAlbum = {
    "uid": uid,
    "data": apiResponse.body
  }
  t.ctx.testChanges.revertAlbums.push(revertAlbum);
}

// this function stores the current information about a photo that will need to be reverted.
export async function helperRevertPhoto (t, uid) {
  const apiResponse = await t.request(`${testcafeconfig.api}photos/${uid}`);
  const revertPhoto = {
    "uid": uid,
    "data": apiResponse.body
  }
  t.ctx.testChanges.revertPhotos.push(revertPhoto);
}

// this function stores the need to remove an album.
export async function helperRemoveAlbum (t, how, id) {
  logMessage(`helperRemoveAlbum ${how} -> ${id}`);
  if (how === "name") {
    const removeAlbum = {
      "name": id,
      "uid": "name"
    }
    t.ctx.testChanges.removeAlbums.push(removeAlbum);
  } else {
    const removeAlbum = {
      "name": "uid",
      "uid": id
    }
    t.ctx.testChanges.removeAlbums.push(removeAlbum);
  }
}

// this function stores the need to remove a label from a photo.
export async function helperRemoveLabelFromPhotos (t, labelUid, photoUid) {
  if (labelUid.length > 0) {
    if (!isNaN(+labelUid)) {
      const removeLabelFromPhoto = {
        "labelUid": labelUid,
        "photoUid" : photoUid
      }
      t.ctx.testChanges.removeLabelFromPhotos.push(removeLabelFromPhoto);
    } else {
      const apiResponse = await t.request(`${testcafeconfig.api}labels?count=1&q=uid:${labelUid}`);
      if (apiResponse.status == 200) {
        const removeLabelFromPhoto = {
          "labelUid": apiResponse.body[0].ID,
          "photoUid" : photoUid
        }
        t.ctx.testChanges.removeLabelFromPhotos.push(removeLabelFromPhoto);
      }
    }
  }
}

// this function stores the need to remove a label from ALL photos.
export async function helperRemoveLabel (t, name) {
  const removeLabel = {
    "name": name
  }
  t.ctx.testChanges.removeLabels.push(removeLabel);
}


// This function will undo what the test has done (to the best of it's ability)
export async function helperAfterEach(t) {
  logMessage("helperAfterEach Queued Requests");
  logMessage(JSON.stringify(t.ctx.testChanges));
  // Revert Albums state
  for (const revertAlbum of t.ctx.testChanges.revertAlbums) {
    const apiResponse = await t.request({
      url: `${testcafeconfig.api}albums/${revertAlbum.uid}`,
      method: 'put',
      body: revertAlbum.data
    });
    // ToDo: handle a bad apiResponse
    if (apiResponse.status != 200 || apiResponse.status === null) { // Ignore Ok
      logMessage("helperAfterEach revert albums " + JSON.stringify(apiResponse));
    }
  }
  
  // Revert Photos state
  // this can not fully restore labels to a photo.
  // if the label has been fully removed, and it matches a keyword, then it will 
  // be restored as a keyword based label.  Otherwise it will be a manual
  // style label.
  for (const revertPhoto of t.ctx.testChanges.revertPhotos) {
    // Revert the photo
    var apiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${revertPhoto.uid}`,
      method: 'put',
      body: revertPhoto.data
    });
    // ToDo: handle a bad apiResponse
    if (apiResponse.status != 200 || apiResponse.status === null) { // Ignore Ok
      logMessage("helperAfterEach revert photo " + JSON.stringify(apiResponse));
    }

    // Loop through the labels in revertPhoto.data and apiResponse.body to add/remove as needed.
    // Remove
    for (const label of apiResponse.body.Labels) {
      const exists = revertPhoto.data.Labels.some(slug => slug.Label.Slug === label.Label.Slug);
      if (!exists) {
        helperRemoveLabelFromPhotos(t, label.Label.ID, revertPhoto.uid);
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
              "Uncertanty": label.Label.Uncertanty
          }
        });
        if (labelApiResponse.status != 200 || labelApiResponse.status === null) { // Ignore Ok
          logMessage("helperAfterEach add label " + JSON.stringify(labelApiResponse));
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
        if (albumApiResponse.status != 200 || albumApiResponse.status === null) { // Ignore Ok
          logMessage("helperAfterEach delete from album " + JSON.stringify(albumApiResponse));
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
        if (albumApiResponse.status != 200 || albumApiResponse.status === null) { // Ignore Ok
          logMessage("helperAfterEach add to album " + JSON.stringify(albumApiResponse));
        }
      }
    }

    // Loop through the files and markers to update as required
    // Invalidate any that shouldn't be there.
    for (const file of apiResponse.body.Files) {
      for (const marker of file.Markers) {
        const rFile = revertPhoto.data.Files.find(fileI => fileI.UID === file.UID)
        if (rFile) {
          const rMarker = rFile.Markers.find(m => m.UID === marker.UID && m.FileUID === marker.FileUID);
          logMessage("helperAfterEach rMarker (1)" + JSON.stringify(rMarker));
          var markerApiResponse;
          if (rMarker) {
            // reset
            logMessage(`reset ${testcafeconfig.api}markers/${rMarker.UID}`);
            markerApiResponse = await t.request({
              url: `${testcafeconfig.api}markers/${rMarker.UID}`,
              method: 'put',
              body: rMarker
            });
          } else {
            // inactivate
            logMessage(`invalidate ${testcafeconfig.api}markers/${marker.UID}`);
            markerApiResponse = await t.request({
              url: `${testcafeconfig.api}markers/${marker.UID}`,
              method: 'put',
              body: {
                "Invalid":true
              }
            });
          }
          if (markerApiResponse.status != 200 || markerApiResponse.status === null) { // Ignore Ok
            logMessage("helperAfterEach sync markers (1) file " + marker.FileUID + " marker " + marker.UID + " " + JSON.stringify(markerApiResponse));
          }

        }
      }
    }
    for (const file of revertPhoto.data.Files) {
      for (const marker of file.Markers) {
        const rMarker = apiResponse.body.Files.find(file => file.Markers.UID === marker.UID && file.Markers.FileUID === marker.FileUID);
        logMessage("helperAfterEach rMarker (2)" + JSON.stringify(rMarker));
        logMessage(`reset ${testcafeconfig.api}markers/${marker.UID}`);
        const markerApiResponse = await t.request({
            url: `${testcafeconfig.api}markers/${marker.UID}`,
            method: 'put',
            body: marker
          });
        if (markerApiResponse.status != 200 || markerApiResponse.status === null) { // Ignore Ok
          logMessage("helperAfterEach sync markers (1) " + JSON.stringify(markerApiResponse));
        }
      }
    }

    // Do the photo again to try the Title again.
    apiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${revertPhoto.uid}`,
      method: 'put',
      body: revertPhoto.data
    });
    // ToDo: handle a bad apiResponse
    if (apiResponse.status != 200 || apiResponse.status === null) { // Ignore Ok
      logMessage("helperAfterEach revert photo again " + JSON.stringify(apiResponse));
    }

  }
  // Remove albums
  for (const removeAlbum of t.ctx.testChanges.removeAlbums) {
    var listApiResponse
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
    if (listApiResponse.status != 200 || listApiResponse.status === null) { // Ignore Ok
      logMessage("helperAfterEach list albums " + JSON.stringify(listApiResponse));
    }
    logMessage("remove album " + JSON.stringify(listApiResponse));
    for (const album of listApiResponse.body) {
      const apiResponse = await t.request({
        url: `${testcafeconfig.api}albums/${album.UID}`,
        method: 'delete',
        params: {
            force: true
        }
      });
      // ToDo: handle a bad apiResponse
      if (apiResponse.status != 200 || apiResponse.status === null && apiResponse.status != 404) { // Ignore Ok and not found
        logMessage("helperAfterEach delete album " + JSON.stringify(apiResponse));
      }
    }
  }
  // Remove Labels from Photos
  for (const removeLabelFromPhoto of t.ctx.testChanges.removeLabelFromPhotos) {
    const apiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${removeLabelFromPhoto.photoUid}/label/${removeLabelFromPhoto.labelUid}`,
      method: 'delete'
    });
    // ToDo: handle a bad apiResponse
    if ((apiResponse.status != 200 && apiResponse.status != 404) || apiResponse.status === null ) { // Ignore Ok and not found
      logMessage("helperAfterEach remove label from photo " + JSON.stringify(apiResponse));
    }
  }
  // Remove Labels
  if (t.ctx.testChanges.removeLabels.length > 0) {
    var labels = [];
    for (const removeLabel of t.ctx.testChanges.removeLabels) {
      const listApiResponse = await t.request({
        url: `${testcafeconfig.api}labels`,
        method: 'get',
        params: {
          count: 10,
          q: `${removeLabel.name}`
        }
      });
      if (listApiResponse.status != 200 || listApiResponse.status === null) { // Ignore Ok
        logMessage("helperAfterEach get labels " + JSON.stringify(listApiResponse));
      }
      for (const label of listApiResponse.body) {
        labels.push(label.UID);
      }
    }
    const apiResponse = await t.request({
      url: `${testcafeconfig.api}batch/labels/delete`,
      method: 'post',
      body: {
        "labels": labels
      }
    });
    if (apiResponse.status != 200 || apiResponse.status === null) { // Ignore Ok
      logMessage("helperAfterEach delete labels " + JSON.stringify(apiResponse));
    }
  }
}