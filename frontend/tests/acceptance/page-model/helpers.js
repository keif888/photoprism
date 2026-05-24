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
  logMessage("albumUid = " + uid);
}

// this function stores the current information about a photo that will need to be reverted.
export async function helperRevertPhoto (t, uid) {
  const apiResponse = await t.request(`${testcafeconfig.api}photos/${uid}`);
  const revertPhoto = {
    "uid": uid,
    "data": apiResponse.body
  }
  t.ctx.testChanges.revertPhotos.push(revertPhoto);
  logMessage("photoUid = " + uid);
}

// this function stores the need to remove an album.
export async function helperRemoveAlbum (t, name) {
  const removeAlbum = {
    "name": name
  }
  t.ctx.testChanges.removeAlbums.push(removeAlbum);
}

// this function stores the need to remove a label from a photo.
export async function helperRemoveLabelFromPhotos (t, labelUid, photoUid) {
  const removeLabelFromPhoto = {
    "labelUid": labelUid,
    "photoUid" : photoUid
  }
  t.ctx.testChanges.removeLabelFromPhotos.push(removeLabelFromPhoto);
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
  // Revert Albums state
  for (const revertAlbum of t.ctx.testChanges.revertAlbums) {
    const apiResponse = await t.request({
      url: `${testcafeconfig.api}albums/${revertAlbum.uid}`,
      method: 'put',
      body: revertAlbum.data
    });
    // ToDo: handle a bad apiResponse
    if (apiResponse.status != 200) { // Ignore Ok
      logMessage(JSON.stringify(apiResponse));
    }
  }
  // Revert Photos state
  for (const revertPhoto of t.ctx.testChanges.revertPhotos) {
    const apiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${revertPhoto.uid}`,
      method: 'put',
      body: revertPhoto.data
    });
    // ToDo: handle a bad apiResponse
    if (apiResponse.status != 200) { // Ignore Ok
      logMessage(JSON.stringify(apiResponse));
    }
  }
  // Remove albums
  for (const removeAlbum of t.ctx.testChanges.removeAlbums) {
    const listApiResponse = await t.request({
      url: `${testcafeconfig.api}albums`,
      method: 'get',
      params: {
        count: 10,
        q: `${removeAlbum.name}`
      }
    });
    if (listApiResponse.status != 200) { // Ignore Ok
      logMessage(JSON.stringify(listApiResponse));
    }
    for (const album of listApiResponse.body) {
      const apiResponse = await t.request({
        url: `${testcafeconfig.api}albums/${album.UID}`,
        method: 'delete',
        params: {
            force: true
        }
      });
      // ToDo: handle a bad apiResponse
      if (apiResponse.status != 200 && apiResponse.status != 404) { // Ignore Ok and not found
        logMessage(JSON.stringify(apiResponse));
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
    if (apiResponse.status != 200 && apiResponse.status != 404) { // Ignore Ok and not found
      logMessage(JSON.stringify(apiResponse));
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
      if (listApiResponse.status != 200) { // Ignore Ok
        logMessage(JSON.stringify(listApiResponse));
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
    if (apiResponse.status != 200) { // Ignore Ok
      logMessage(JSON.stringify(apiResponse));
    }
  }
}