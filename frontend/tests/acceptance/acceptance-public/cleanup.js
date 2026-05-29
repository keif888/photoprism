import { Selector } from "testcafe";
import testcafeconfig from "../../testcafeconfig.json";
import { helperBeforeEach, helperAfterEach, helperRemoveAlbum, helperRevertAlbum, helperRevertPhoto, helperRemoveLabelFromPhotos, helperRemoveLabel, logMessage } from "../page-model/helpers";

fixture`Test helper`
.page`${testcafeconfig.url}`;

function cleanAlbumsFilesAndLabels(jsonBody) {
  let items = JSON.parse(JSON.stringify(jsonBody.Albums));
  jsonBody.Albums.length = 0;
  for (let item of items) {
    delete item.UpdatedAt;
    jsonBody.Albums.push(item);
  }
  items = JSON.parse(JSON.stringify(jsonBody.Files));
  jsonBody.Files.length = 0;
  for (let item of items) {
    delete item.UpdatedAt;
    jsonBody.Files.push(item);
  }
  items = JSON.parse(JSON.stringify(jsonBody.Labels));
  jsonBody.Labels.length = 0;
  for (let item of items) {
    delete item.Label.UpdatedAt;
    delete item.LabelSrc; // can't be set by API
    delete item.Uncertainty; // when Src changes, the uncertainty changes.
    jsonBody.Labels.push(item);
  }
}

test.meta("testID", "cleanup-001").meta({ type: "short", mode: "api" })("Common: Cleanup helperRevertAlbum remove album cover photo and change caption/description of existing album", async (t) => {
    await helperBeforeEach(t);
    let beforeAlbumResponse = await t.request({
        url: `${testcafeconfig.api}albums`,
        method: 'get',
        params: {
          count: 1,
          q: `Christmas type:album`
        }
      });
    await t.expect(beforeAlbumResponse.status).eql(200);
    const albumUID = beforeAlbumResponse.body[0].UID;
    beforeAlbumResponse = await t.request(`${testcafeconfig.api}albums/${albumUID}`);
    await t.expect(beforeAlbumResponse.status).eql(200);
    await helperRevertAlbum(t, albumUID);
    // Change the name on the album
    let apiResponse = await t.request({
      url: `${testcafeconfig.api}albums/${albumUID}`,
      method: 'put',
      body: {
        "Caption": "Cleanup test data",
        "Description": "This should be removed"
      }
    });
    await t.expect(apiResponse.status).eql(200);

    // remove pqmxlr7188hz4bih from album "Christmas", which is the album cover photo
    apiResponse = await t.request({
      url: `${testcafeconfig.api}albums/${albumUID}/photos`,
      method: 'delete',
      body: {
        "photos": [ "pqmxlr7188hz4bih" ]  // This should be ok.
      }
    });
    await t.expect(apiResponse.status).eql(200);

    await helperAfterEach(t);
    const afterAlbumResponse = await t.request(`${testcafeconfig.api}albums/${albumUID}`);
    await t.expect(afterAlbumResponse).notEql(beforeAlbumResponse);
    // Remove the fields that are impacted by changes
    delete beforeAlbumResponse.body.UpdatedAt; // Will change
    delete afterAlbumResponse.body.UpdatedAt;
    delete beforeAlbumResponse.body.ThumbSrc; // Will change on 1st run, wont on subsequent
    delete afterAlbumResponse.body.ThumbSrc;
    delete beforeAlbumResponse.headers["content-length"]; // Will change (timestamp)
    delete afterAlbumResponse.headers["content-length"];
    delete beforeAlbumResponse.headers.date; // May change if second ticks over
    delete afterAlbumResponse.headers.date;
    await t.expect(afterAlbumResponse).eql(beforeAlbumResponse);
});

test.meta("testID", "cleanup-002").meta({ type: "short", mode: "api" })("Common: Cleanup helperRevertAlbum remove Garden album", async (t) => {
    await helperBeforeEach(t);
    const albumUID = "arkgush1tdwk4fsy";
    let beforeAlbumResponse = await t.request(`${testcafeconfig.api}albums/${albumUID}`);
    await t.expect(beforeAlbumResponse.status).eql(200);
    await helperRevertAlbum(t, albumUID);
    // Soft Delete the album
    let apiResponse = await t.request({
      url: `${testcafeconfig.api}albums/${albumUID}`,
      method: 'delete',
      params: {
        force: false
      }      
    });
    await t.expect(apiResponse.status).eql(200);

    await helperAfterEach(t);
    let afterAlbumResponse = await t.request(`${testcafeconfig.api}albums/${albumUID}`);
    await t.expect(afterAlbumResponse).notEql(beforeAlbumResponse);
    await t.expect(afterAlbumResponse.status).eql(200); // A deleted album with the SAME CreatedBy as the current user will be undeleted
    // Remove the fields that are impacted by changes
    delete beforeAlbumResponse.body.UpdatedAt;
    delete afterAlbumResponse.body.UpdatedAt;
    delete beforeAlbumResponse.body.ThumbSrc;
    delete afterAlbumResponse.body.ThumbSrc;
    delete beforeAlbumResponse.headers["content-length"];
    delete afterAlbumResponse.headers["content-length"];
    delete beforeAlbumResponse.headers.date;
    delete afterAlbumResponse.headers.date;
    await t.expect(afterAlbumResponse).eql(beforeAlbumResponse);
});

test.meta("testID", "cleanup-003").meta({ type: "short", mode: "api" })("Common: Cleanup helperRevertAlbum remove Holiday album", async (t) => {
    await helperBeforeEach(t);
    let beforeAlbumResponse = await t.request({
        url: `${testcafeconfig.api}albums`,
        method: 'get',
        params: {
          count: 1,
          q: `Holiday type:album`
        }
      });
    await t.expect(beforeAlbumResponse.status).eql(200);
    const albumUID = beforeAlbumResponse.body[0].UID;
    await t.expect(albumUID).eql("aqmxlt22ilujuxux", "Test requires Holiday album to not have been deleted already");
    beforeAlbumResponse = await t.request(`${testcafeconfig.api}albums/${albumUID}`);
    await t.expect(beforeAlbumResponse.status).eql(200);
    await helperRevertAlbum(t, albumUID);
    // Soft Delete the album
    let apiResponse = await t.request({
      url: `${testcafeconfig.api}albums/${albumUID}`,
      method: 'delete',
      params: {
        force: false
      }      
    });
    await t.expect(apiResponse.status).eql(200);

    await helperAfterEach(t);
    let afterAlbumResponse = await t.request(`${testcafeconfig.api}albums/${albumUID}`);
    await t.expect(afterAlbumResponse).notEql(beforeAlbumResponse);
    await t.expect(afterAlbumResponse.status).eql(404); // A deleted album with a null or different CreatedBy will change UID
    afterAlbumResponse = await t.request({
      url: `${testcafeconfig.api}albums`,
      method: 'get',
      params: {
        count: 1,
        q: `Holiday type:album`
      }
    });
    await t.expect(afterAlbumResponse.status).eql(200);
    afterAlbumResponse = await t.request(`${testcafeconfig.api}albums/${afterAlbumResponse.body[0].UID}`);
    await t.expect(afterAlbumResponse.status).eql(200);
    // Remove the fields that are impacted by changes
    delete beforeAlbumResponse.body.UID;
    delete afterAlbumResponse.body.UID;
    delete beforeAlbumResponse.body.ID;
    delete afterAlbumResponse.body.ID;
    delete beforeAlbumResponse.body.UpdatedAt;
    delete afterAlbumResponse.body.UpdatedAt;
    delete beforeAlbumResponse.body.CreatedAt;
    delete afterAlbumResponse.body.CreatedAt;
    delete beforeAlbumResponse.body.CreatedBy;
    delete afterAlbumResponse.body.CreatedBy;
    delete beforeAlbumResponse.body.ThumbSrc;
    delete afterAlbumResponse.body.ThumbSrc;
    delete beforeAlbumResponse.headers["content-length"];
    delete afterAlbumResponse.headers["content-length"];
    delete beforeAlbumResponse.headers.date;
    delete afterAlbumResponse.headers.date;
    await t.expect(afterAlbumResponse).eql(beforeAlbumResponse);
});

test.meta("testID", "cleanup-004").meta({ type: "short", mode: "api" })("Common: Cleanup helperRevertPhoto revert titles and details", async (t) => {
    await helperBeforeEach(t);
    let beforePhotoResponse = await t.request({
        url: `${testcafeconfig.api}photos`,
        method: 'get',
        params: {
          count: 1,
          q: `geo:true camera:apple`
        }
      });
    await t.expect(beforePhotoResponse.status).eql(200);
    const photoUID = beforePhotoResponse.body[0].UID;
    beforePhotoResponse = await t.request(`${testcafeconfig.api}photos/${photoUID}`);
    await t.expect(beforePhotoResponse.status).eql(200);
    await helperRevertPhoto(t, photoUID);
    // Change the name and other stuff on the photo
    let apiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${photoUID}`,
      method: 'put',
      body: {
        "Title": "Cleanup test data",
        "Description": "This should be removed",
        "CameraID": 7,
        "LensID": 10,
        "CellID": "s2:47a85a634bcc",
        "PlaceID": "de:ukLS8nroIoB7"
      }
    });
    await t.expect(apiResponse.status).eql(200);

    await helperAfterEach(t);

    let afterPhotoResponse = await t.request(`${testcafeconfig.api}photos/${photoUID}`);
    await t.expect(afterPhotoResponse).notEql(beforePhotoResponse);
    // Remove the fields that are impacted by changes
    delete beforePhotoResponse.headers["content-length"]; // Will change (timestamp)
    delete afterPhotoResponse.headers["content-length"];
    delete beforePhotoResponse.headers.date; // May change if second ticks over
    delete afterPhotoResponse.headers.date;
    delete beforePhotoResponse.body.UpdatedAt;
    delete afterPhotoResponse.body.UpdatedAt;
    delete beforePhotoResponse.body.EditedAt;
    delete afterPhotoResponse.body.EditedAt;
    delete beforePhotoResponse.body.Details.UpdatedAt;
    delete afterPhotoResponse.body.Details.UpdatedAt;
    cleanAlbumsFilesAndLabels(beforePhotoResponse.body);
    cleanAlbumsFilesAndLabels(afterPhotoResponse.body);
    await t.expect(afterPhotoResponse).eql(beforePhotoResponse);
})

test.meta("testID", "cleanup-005").meta({ type: "short", mode: "api" })("Common: Cleanup helperRevertPhoto revert Labels", async (t) => {
    await helperBeforeEach(t);
    const stamp = Date.now();
    const labelTitle = `SidebarEditLabel-${stamp}`;
    let beforePhotoResponse = await t.request({
        url: `${testcafeconfig.api}photos`,
        method: 'get',
        params: {
          count: 1,
          q: `label:cat`
        }
      });
    await t.expect(beforePhotoResponse.status).eql(200);
    const photoUID = beforePhotoResponse.body[0].UID;
    beforePhotoResponse = await t.request(`${testcafeconfig.api}photos/${photoUID}`);
    await t.expect(beforePhotoResponse.status).eql(200);
    const labelID = 11; // The Cat label! beforePhotoResponse.body.Labels[0].ID;
    await helperRevertPhoto(t, photoUID);
    // Remove a label from the photo
    let apiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${photoUID}/label/${labelID}`,
      method: 'delete'
    });
    await t.expect(apiResponse.status).eql(200);

    // Add a manual label
    const labelApiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${photoUID}/label`,
      method: 'post',
      body: {
          "Description": "Testing Label",
          "Favorite": false,
          "Name": labelTitle,
          "Priority": 0,
          "Uncertainty": 0
      }
    });
    await t.expect(labelApiResponse.status).eql(200);

    await helperAfterEach(t);

    const afterPhotoResponse = await t.request(`${testcafeconfig.api}photos/${photoUID}`);
    await t.expect(afterPhotoResponse).notEql(beforePhotoResponse);
    // Remove the fields that are impacted by changes
    delete beforePhotoResponse.headers["content-length"]; // Will change (timestamp)
    delete afterPhotoResponse.headers["content-length"];
    delete beforePhotoResponse.headers.date; // May change if second ticks over
    delete afterPhotoResponse.headers.date;
    if (!beforePhotoResponse.body.Title.includes("Cat")) {
      // If the title didn't include cat, then the label change to Manual 100% certainty MAY add cat to the Title.
      delete beforePhotoResponse.body.Title;
      delete afterPhotoResponse.body.Title;
    }
    delete beforePhotoResponse.body.UpdatedAt;
    delete afterPhotoResponse.body.UpdatedAt;
    delete beforePhotoResponse.body.EditedAt;
    delete afterPhotoResponse.body.EditedAt;
    delete beforePhotoResponse.body.Details.UpdatedAt;
    delete afterPhotoResponse.body.Details.UpdatedAt;
    cleanAlbumsFilesAndLabels(beforePhotoResponse.body);
    cleanAlbumsFilesAndLabels(afterPhotoResponse.body);
    await t.expect(afterPhotoResponse).eql(beforePhotoResponse);
})

// This test will leave a junk label behind if it fails, as it's testing that manual labels are reverted correctly.
test.meta("testID", "cleanup-006").meta({ type: "short", mode: "api" })("Common: Cleanup helperRevertPhoto revert manual deleted Label", async (t) => {
    await helperBeforeEach(t);
    const stamp = Date.now();
    const labelTitle = `SidebarEditLabel-${stamp}`;
    let beforePhotoResponse = await t.request({
        url: `${testcafeconfig.api}photos`,
        method: 'get',
        params: {
          count: 1
        }
      });
    await t.expect(beforePhotoResponse.status).eql(200);
    const photoUID = beforePhotoResponse.body[0].UID;
    // Add a manual label
    let labelApiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${photoUID}/label`,
      method: 'post',
      body: {
          "Description": "Testing Label",
          "Favorite": false,
          "Name": labelTitle,
          "Priority": 0,
          "Uncertainty": 0
      }
    });
    await t.expect(labelApiResponse.status).eql(200);
    labelApiResponse = await t.request({
      url: `${testcafeconfig.api}labels`,
      method: 'get',
      params: {
        count: 1,
        q: labelTitle
      }
    });
    await t.expect(labelApiResponse.status).eql(200);
    const labelID = labelApiResponse.body[0].ID;

    beforePhotoResponse = await t.request(`${testcafeconfig.api}photos/${photoUID}`);
    await t.expect(beforePhotoResponse.status).eql(200);
    await helperRevertPhoto(t, photoUID);

    // Remove the manual label from the photo
    let apiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${photoUID}/label/${labelID}`,
      method: 'delete'
    });
    await t.expect(apiResponse.status).eql(200);
    await helperAfterEach(t);

    const afterPhotoResponse = await t.request(`${testcafeconfig.api}photos/${photoUID}`);
    await t.expect(afterPhotoResponse).notEql(beforePhotoResponse);
    // Remove the fields that are impacted by changes
    delete beforePhotoResponse.headers["content-length"]; // Will change (timestamp)
    delete afterPhotoResponse.headers["content-length"];
    delete beforePhotoResponse.headers.date; // May change if second ticks over
    delete afterPhotoResponse.headers.date;
    delete beforePhotoResponse.body.UpdatedAt;
    delete afterPhotoResponse.body.UpdatedAt;
    delete beforePhotoResponse.body.EditedAt;
    delete afterPhotoResponse.body.EditedAt;
    delete beforePhotoResponse.body.Details.UpdatedAt;
    delete afterPhotoResponse.body.Details.UpdatedAt;
    cleanAlbumsFilesAndLabels(beforePhotoResponse.body);
    cleanAlbumsFilesAndLabels(afterPhotoResponse.body);
    await t.expect(afterPhotoResponse).eql(beforePhotoResponse);

    // Remove the manual label from the photo
    apiResponse = await t.request({
      url: `${testcafeconfig.api}photos/${photoUID}/label/${labelID}`,
      method: 'delete'
    });
    await t.expect(apiResponse.status).eql(200);
})

// This test will leave a junk album behind if it fails, as it's testing that manual albums are reverted correctly from a photo.
test.meta("testID", "cleanup-007").meta({ type: "short", mode: "api" })("Common: Cleanup helperRevertPhoto revert Albums", async (t) => {
    await helperBeforeEach(t);
    const stamp = Date.now();
    const albumTitle = `SidebarEditAlbum-${stamp}`;
    let beforePhotoResponse = await t.request({
        url: `${testcafeconfig.api}photos`,
        method: 'get',
        params: {
          count: 1,
          q: `album:garden private:false`
        }
      });
    await t.expect(beforePhotoResponse.status).eql(200);
    const photoUID = beforePhotoResponse.body[0].UID;
    beforePhotoResponse = await t.request(`${testcafeconfig.api}photos/${photoUID}`);
    await t.expect(beforePhotoResponse.status).eql(200);

    await helperRevertPhoto(t, photoUID);

    const albumUID = "arkgush1tdwk4fsy";
    // Remove an existing album from the photo
    let apiResponse = await t.request({
      url: `${testcafeconfig.api}albums/${albumUID}/photos`,
      method: 'delete',
      body: {
        "photos": [ photoUID ]
      }
    });
    await t.expect(apiResponse.status).eql(200);

    // Add a manual album
    let albumApiResponse = await t.request({
      url: `${testcafeconfig.api}albums`,
      method: 'post',
      body: {
          "Caption": "cleanup-007",
          "Title": albumTitle
      }
    });
    await t.expect(albumApiResponse.status).eql(201);

    // Add photo to the new album
    const newAlbumUID = albumApiResponse.body.UID;

    albumApiResponse = await t.request({
      url: `${testcafeconfig.api}albums/${newAlbumUID}/photos`,
      method: 'post',
      body: {
        "photos": [ photoUID ]
      }
    });
    await t.expect(albumApiResponse.status).eql(200);
    logMessage(JSON.stringify(albumApiResponse));

    await helperAfterEach(t);

    const afterPhotoResponse = await t.request(`${testcafeconfig.api}photos/${photoUID}`);
    await t.expect(afterPhotoResponse).notEql(beforePhotoResponse);
    // Remove the fields that are impacted by changes
    delete beforePhotoResponse.headers["content-length"]; // Will change (timestamp)
    delete afterPhotoResponse.headers["content-length"];
    delete beforePhotoResponse.headers.date; // May change if second ticks over
    delete afterPhotoResponse.headers.date;
    delete beforePhotoResponse.body.UpdatedAt;
    delete afterPhotoResponse.body.UpdatedAt;
    delete beforePhotoResponse.body.EditedAt;
    delete afterPhotoResponse.body.EditedAt;
    delete beforePhotoResponse.body.Details.UpdatedAt;
    delete afterPhotoResponse.body.Details.UpdatedAt;
    cleanAlbumsFilesAndLabels(beforePhotoResponse.body);
    cleanAlbumsFilesAndLabels(afterPhotoResponse.body);
    await t.expect(afterPhotoResponse).eql(beforePhotoResponse);

    await helperBeforeEach(t);
    await helperRemoveAlbum(t, "name", albumTitle);
    await helperAfterEach(t);

    albumApiResponse = await t.request({
      url: `${testcafeconfig.api}albums/${newAlbumUID}`,
      method: 'get'
    });
    await t.expect(albumApiResponse.status).eql(404);
})
