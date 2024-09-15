import { Selector } from "testcafe";
import testcafeconfig from "../../testcafeconfig.json";
import { ClientFunction } from "testcafe";
import Menu from "../page-model/menu";
import Toolbar from "../page-model/toolbar";
import ContextMenu from "../page-model/context-menu";
import Photo from "../page-model/photo";
import PhotoViewer from "../page-model/photoviewer";
import Page from "../page-model/page";
import PhotoEdit from "../page-model/photo-edit";

const scroll = ClientFunction((x, y) => window.scrollTo(x, y));
const getcurrentPosition = ClientFunction(() => window.pageYOffset);

fixture`Test photos Geo64bit`.page`${testcafeconfig.url}`;

const menu = new Menu();
const toolbar = new Toolbar();
const contextmenu = new ContextMenu();
const photo = new Photo();
const photoviewer = new PhotoViewer();
const page = new Page();
const photoedit = new PhotoEdit();


test.meta("testID", "photosgeo64-001").meta({ type: "short", mode: "public" })(
  "Common: Add 2 decimal Location",
  async (t) => {
    await menu.openPage("browse");
    await toolbar.setFilter("view", "Cards");
    const FirstPhotoUid = await photo.getNthPhotoUid("image", 0);
    await t.click(page.cardTitle.withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.latitude.visible).ok();

    await t.click(photoedit.dialogNext);

    await t.expect(photoedit.dialogPrevious.getAttribute("disabled")).notEql("disabled");

    await t.click(photoedit.dialogPrevious).click(photoedit.dialogClose);
    await photoviewer.openPhotoViewer("uid", FirstPhotoUid);
    await photoviewer.triggerPhotoViewerAction("edit");
    const FirstPhotoTitle = await photoedit.title.value;
    const FirstPhotoLocalTime = await photoedit.localTime.value;
    const FirstPhotoDay = await photoedit.day.value;
    const FirstPhotoMonth = await photoedit.month.value;
    const FirstPhotoYear = await photoedit.year.value;
    const FirstPhotoTimezone = await photoedit.timezone.value;
    const FirstPhotoLatitude = await photoedit.latitude.value;
    const FirstPhotoLongitude = await photoedit.longitude.value;
    const FirstPhotoAltitude = await photoedit.altitude.value;
    const FirstPhotoCountry = await photoedit.country.value;
    const FirstPhotoCamera = await photoedit.camera.innerText;
    const FirstPhotoIso = await photoedit.iso.value;
    const FirstPhotoExposure = await photoedit.exposure.value;
    const FirstPhotoLens = await photoedit.lens.innerText;
    const FirstPhotoFnumber = await photoedit.fnumber.value;
    const FirstPhotoFocalLength = await photoedit.focallength.value;
    const FirstPhotoSubject = await photoedit.subject.value;
    const FirstPhotoArtist = await photoedit.artist.value;
    const FirstPhotoCopyright = await photoedit.copyright.value;
    const FirstPhotoLicense = await photoedit.license.value;
    const FirstPhotoDescription = await photoedit.description.value;
    const FirstPhotoKeywords = await photoedit.keywords.value;
    const FirstPhotoNotes = await photoedit.notes.value;

    await t
      .typeText(photoedit.title, "Not saved photo title", { replace: true })
      .click(photoedit.detailsClose)
      .click(Selector("button.action-title-edit").withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.title.value).eql(FirstPhotoTitle);

    await photoedit.editPhoto(
      "New Photo Title",
      "Europe/Moscow",
      "15",
      "07",
      "2019",
      "04:30:30",
      "-1",
      "41.15",
      "20.16",
      "32",
      "1/32",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      ", cat, love",
      "Some notes"
    );
    if (t.browser.platform === "mobile") {
      await t.eval(() => location.reload());
    } else {
      await toolbar.triggerToolbarAction("reload");
    }
    await toolbar.search("uid:" + FirstPhotoUid);

    await t
      .expect(page.cardTitle.withAttribute("data-uid", FirstPhotoUid).innerText)
      .eql("New Photo Title");

    await photo.triggerHoverAction("uid", FirstPhotoUid, "select");
    await contextmenu.triggerContextMenuAction("edit", "");

    //const expectedValues = [{ FirstPhotoTitle: photoedit.title }, { "bluh bla": photoedit.day }];
    /*const expectedValues = [
    [FirstPhotoTitle, photoedit.title],
    ["blah", photoedit.day],
  ];
  await photoedit.checkEditFormValuesNewNew(expectedValues);*/

    await photoedit.checkEditFormValues(
      "New Photo Title",
      "15",
      "07",
      "2019",
      "04:30:30",
      "Europe/Moscow",
      "Albania",
      "-1",
      "41.15",
      "20.16",
      "",
      "32",
      "1/32",
      "",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      "cat",
      "Some notes"
    );

    await photoedit.undoPhotoEdit(
      FirstPhotoTitle,
      FirstPhotoTimezone,
      FirstPhotoDay,
      FirstPhotoMonth,
      FirstPhotoYear,
      FirstPhotoLocalTime,
      FirstPhotoAltitude,
      FirstPhotoLatitude,
      FirstPhotoLongitude,
      FirstPhotoCountry,
      FirstPhotoIso,
      FirstPhotoExposure,
      FirstPhotoFnumber,
      FirstPhotoFocalLength,
      FirstPhotoSubject,
      FirstPhotoArtist,
      FirstPhotoCopyright,
      FirstPhotoLicense,
      FirstPhotoDescription,
      FirstPhotoKeywords,
      FirstPhotoNotes
    );
    await contextmenu.checkContextMenuCount("1");
    await contextmenu.clearSelection();
  }
);

test.meta("testID", "photosgeo64-002").meta({ type: "short", mode: "public" })(
  "Common: Add 3 decimal Location",
  async (t) => {
    await menu.openPage("browse");
    await toolbar.setFilter("view", "Cards");
    const FirstPhotoUid = await photo.getNthPhotoUid("image", 0);
    await t.click(page.cardTitle.withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.latitude.visible).ok();

    await t.click(photoedit.dialogNext);

    await t.expect(photoedit.dialogPrevious.getAttribute("disabled")).notEql("disabled");

    await t.click(photoedit.dialogPrevious).click(photoedit.dialogClose);
    await photoviewer.openPhotoViewer("uid", FirstPhotoUid);
    await photoviewer.triggerPhotoViewerAction("edit");
    const FirstPhotoTitle = await photoedit.title.value;
    const FirstPhotoLocalTime = await photoedit.localTime.value;
    const FirstPhotoDay = await photoedit.day.value;
    const FirstPhotoMonth = await photoedit.month.value;
    const FirstPhotoYear = await photoedit.year.value;
    const FirstPhotoTimezone = await photoedit.timezone.value;
    const FirstPhotoLatitude = await photoedit.latitude.value;
    const FirstPhotoLongitude = await photoedit.longitude.value;
    const FirstPhotoAltitude = await photoedit.altitude.value;
    const FirstPhotoCountry = await photoedit.country.value;
    const FirstPhotoCamera = await photoedit.camera.innerText;
    const FirstPhotoIso = await photoedit.iso.value;
    const FirstPhotoExposure = await photoedit.exposure.value;
    const FirstPhotoLens = await photoedit.lens.innerText;
    const FirstPhotoFnumber = await photoedit.fnumber.value;
    const FirstPhotoFocalLength = await photoedit.focallength.value;
    const FirstPhotoSubject = await photoedit.subject.value;
    const FirstPhotoArtist = await photoedit.artist.value;
    const FirstPhotoCopyright = await photoedit.copyright.value;
    const FirstPhotoLicense = await photoedit.license.value;
    const FirstPhotoDescription = await photoedit.description.value;
    const FirstPhotoKeywords = await photoedit.keywords.value;
    const FirstPhotoNotes = await photoedit.notes.value;

    await t
      .typeText(photoedit.title, "Not saved photo title", { replace: true })
      .click(photoedit.detailsClose)
      .click(Selector("button.action-title-edit").withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.title.value).eql(FirstPhotoTitle);

    await photoedit.editPhoto(
      "New Photo Title",
      "Europe/Moscow",
      "15",
      "07",
      "2019",
      "05:30:30",
      "-1",
      "41.153",
      "20.168",
      "32",
      "1/32",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      ", cat, love",
      "Some notes"
    );
    if (t.browser.platform === "mobile") {
      await t.eval(() => location.reload());
    } else {
      await toolbar.triggerToolbarAction("reload");
    }
    await toolbar.search("uid:" + FirstPhotoUid);

    await t
      .expect(page.cardTitle.withAttribute("data-uid", FirstPhotoUid).innerText)
      .eql("New Photo Title");

    await photo.triggerHoverAction("uid", FirstPhotoUid, "select");
    await contextmenu.triggerContextMenuAction("edit", "");

    //const expectedValues = [{ FirstPhotoTitle: photoedit.title }, { "bluh bla": photoedit.day }];
    /*const expectedValues = [
    [FirstPhotoTitle, photoedit.title],
    ["blah", photoedit.day],
  ];
  await photoedit.checkEditFormValuesNewNew(expectedValues);*/

    await photoedit.checkEditFormValues(
      "New Photo Title",
      "15",
      "07",
      "2019",
      "05:30:30",
      "Europe/Moscow",
      "Albania",
      "-1",
      "41.153",
      "20.168",
      "",
      "32",
      "1/32",
      "",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      "cat",
      "Some notes"
    );

    await photoedit.undoPhotoEdit(
      FirstPhotoTitle,
      FirstPhotoTimezone,
      FirstPhotoDay,
      FirstPhotoMonth,
      FirstPhotoYear,
      FirstPhotoLocalTime,
      FirstPhotoAltitude,
      FirstPhotoLatitude,
      FirstPhotoLongitude,
      FirstPhotoCountry,
      FirstPhotoIso,
      FirstPhotoExposure,
      FirstPhotoFnumber,
      FirstPhotoFocalLength,
      FirstPhotoSubject,
      FirstPhotoArtist,
      FirstPhotoCopyright,
      FirstPhotoLicense,
      FirstPhotoDescription,
      FirstPhotoKeywords,
      FirstPhotoNotes
    );
    await contextmenu.checkContextMenuCount("1");
    await contextmenu.clearSelection();
  }
);


test.meta("testID", "photosgeo64-002").meta({ type: "short", mode: "public" })(
  "Common: Add 4 decimal Location",
  async (t) => {
    await menu.openPage("browse");
    await toolbar.setFilter("view", "Cards");
    const FirstPhotoUid = await photo.getNthPhotoUid("image", 0);
    await t.click(page.cardTitle.withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.latitude.visible).ok();

    await t.click(photoedit.dialogNext);

    await t.expect(photoedit.dialogPrevious.getAttribute("disabled")).notEql("disabled");

    await t.click(photoedit.dialogPrevious).click(photoedit.dialogClose);
    await photoviewer.openPhotoViewer("uid", FirstPhotoUid);
    await photoviewer.triggerPhotoViewerAction("edit");
    const FirstPhotoTitle = await photoedit.title.value;
    const FirstPhotoLocalTime = await photoedit.localTime.value;
    const FirstPhotoDay = await photoedit.day.value;
    const FirstPhotoMonth = await photoedit.month.value;
    const FirstPhotoYear = await photoedit.year.value;
    const FirstPhotoTimezone = await photoedit.timezone.value;
    const FirstPhotoLatitude = await photoedit.latitude.value;
    const FirstPhotoLongitude = await photoedit.longitude.value;
    const FirstPhotoAltitude = await photoedit.altitude.value;
    const FirstPhotoCountry = await photoedit.country.value;
    const FirstPhotoCamera = await photoedit.camera.innerText;
    const FirstPhotoIso = await photoedit.iso.value;
    const FirstPhotoExposure = await photoedit.exposure.value;
    const FirstPhotoLens = await photoedit.lens.innerText;
    const FirstPhotoFnumber = await photoedit.fnumber.value;
    const FirstPhotoFocalLength = await photoedit.focallength.value;
    const FirstPhotoSubject = await photoedit.subject.value;
    const FirstPhotoArtist = await photoedit.artist.value;
    const FirstPhotoCopyright = await photoedit.copyright.value;
    const FirstPhotoLicense = await photoedit.license.value;
    const FirstPhotoDescription = await photoedit.description.value;
    const FirstPhotoKeywords = await photoedit.keywords.value;
    const FirstPhotoNotes = await photoedit.notes.value;

    await t
      .typeText(photoedit.title, "Not saved photo title", { replace: true })
      .click(photoedit.detailsClose)
      .click(Selector("button.action-title-edit").withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.title.value).eql(FirstPhotoTitle);

    await photoedit.editPhoto(
      "New Photo Title",
      "Europe/Moscow",
      "15",
      "07",
      "2019",
      "05:30:30",
      "-1",
      "41.1533",
      "20.1683",
      "32",
      "1/32",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      ", cat, love",
      "Some notes"
    );
    if (t.browser.platform === "mobile") {
      await t.eval(() => location.reload());
    } else {
      await toolbar.triggerToolbarAction("reload");
    }
    await toolbar.search("uid:" + FirstPhotoUid);

    await t
      .expect(page.cardTitle.withAttribute("data-uid", FirstPhotoUid).innerText)
      .eql("New Photo Title");

    await photo.triggerHoverAction("uid", FirstPhotoUid, "select");
    await contextmenu.triggerContextMenuAction("edit", "");

    //const expectedValues = [{ FirstPhotoTitle: photoedit.title }, { "bluh bla": photoedit.day }];
    /*const expectedValues = [
    [FirstPhotoTitle, photoedit.title],
    ["blah", photoedit.day],
  ];
  await photoedit.checkEditFormValuesNewNew(expectedValues);*/

    await photoedit.checkEditFormValues(
      "New Photo Title",
      "15",
      "07",
      "2019",
      "05:30:30",
      "Europe/Moscow",
      "Albania",
      "-1",
      "41.1533",
      "20.1683",
      "",
      "32",
      "1/32",
      "",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      "cat",
      "Some notes"
    );

    await photoedit.undoPhotoEdit(
      FirstPhotoTitle,
      FirstPhotoTimezone,
      FirstPhotoDay,
      FirstPhotoMonth,
      FirstPhotoYear,
      FirstPhotoLocalTime,
      FirstPhotoAltitude,
      FirstPhotoLatitude,
      FirstPhotoLongitude,
      FirstPhotoCountry,
      FirstPhotoIso,
      FirstPhotoExposure,
      FirstPhotoFnumber,
      FirstPhotoFocalLength,
      FirstPhotoSubject,
      FirstPhotoArtist,
      FirstPhotoCopyright,
      FirstPhotoLicense,
      FirstPhotoDescription,
      FirstPhotoKeywords,
      FirstPhotoNotes
    );
    await contextmenu.checkContextMenuCount("1");
    await contextmenu.clearSelection();
  }
);

test.meta("testID", "photosgeo64-002").meta({ type: "short", mode: "public" })(
  "Common: Add 5 decimal Location",
  async (t) => {
    await menu.openPage("browse");
    await toolbar.setFilter("view", "Cards");
    const FirstPhotoUid = await photo.getNthPhotoUid("image", 0);
    await t.click(page.cardTitle.withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.latitude.visible).ok();

    await t.click(photoedit.dialogNext);

    await t.expect(photoedit.dialogPrevious.getAttribute("disabled")).notEql("disabled");

    await t.click(photoedit.dialogPrevious).click(photoedit.dialogClose);
    await photoviewer.openPhotoViewer("uid", FirstPhotoUid);
    await photoviewer.triggerPhotoViewerAction("edit");
    const FirstPhotoTitle = await photoedit.title.value;
    const FirstPhotoLocalTime = await photoedit.localTime.value;
    const FirstPhotoDay = await photoedit.day.value;
    const FirstPhotoMonth = await photoedit.month.value;
    const FirstPhotoYear = await photoedit.year.value;
    const FirstPhotoTimezone = await photoedit.timezone.value;
    const FirstPhotoLatitude = await photoedit.latitude.value;
    const FirstPhotoLongitude = await photoedit.longitude.value;
    const FirstPhotoAltitude = await photoedit.altitude.value;
    const FirstPhotoCountry = await photoedit.country.value;
    const FirstPhotoCamera = await photoedit.camera.innerText;
    const FirstPhotoIso = await photoedit.iso.value;
    const FirstPhotoExposure = await photoedit.exposure.value;
    const FirstPhotoLens = await photoedit.lens.innerText;
    const FirstPhotoFnumber = await photoedit.fnumber.value;
    const FirstPhotoFocalLength = await photoedit.focallength.value;
    const FirstPhotoSubject = await photoedit.subject.value;
    const FirstPhotoArtist = await photoedit.artist.value;
    const FirstPhotoCopyright = await photoedit.copyright.value;
    const FirstPhotoLicense = await photoedit.license.value;
    const FirstPhotoDescription = await photoedit.description.value;
    const FirstPhotoKeywords = await photoedit.keywords.value;
    const FirstPhotoNotes = await photoedit.notes.value;

    await t
      .typeText(photoedit.title, "Not saved photo title", { replace: true })
      .click(photoedit.detailsClose)
      .click(Selector("button.action-title-edit").withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.title.value).eql(FirstPhotoTitle);

    await photoedit.editPhoto(
      "New Photo Title",
      "Europe/Moscow",
      "15",
      "07",
      "2019",
      "05:30:30",
      "-1",
      "41.15333",
      "20.16833",
      "32",
      "1/32",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      ", cat, love",
      "Some notes"
    );
    if (t.browser.platform === "mobile") {
      await t.eval(() => location.reload());
    } else {
      await toolbar.triggerToolbarAction("reload");
    }
    await toolbar.search("uid:" + FirstPhotoUid);

    await t
      .expect(page.cardTitle.withAttribute("data-uid", FirstPhotoUid).innerText)
      .eql("New Photo Title");

    await photo.triggerHoverAction("uid", FirstPhotoUid, "select");
    await contextmenu.triggerContextMenuAction("edit", "");

    //const expectedValues = [{ FirstPhotoTitle: photoedit.title }, { "bluh bla": photoedit.day }];
    /*const expectedValues = [
    [FirstPhotoTitle, photoedit.title],
    ["blah", photoedit.day],
  ];
  await photoedit.checkEditFormValuesNewNew(expectedValues);*/

    await photoedit.checkEditFormValues(
      "New Photo Title",
      "15",
      "07",
      "2019",
      "05:30:30",
      "Europe/Moscow",
      "Albania",
      "-1",
      "41.15333",
      "20.16833",
      "",
      "32",
      "1/32",
      "",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      "cat",
      "Some notes"
    );

    await photoedit.undoPhotoEdit(
      FirstPhotoTitle,
      FirstPhotoTimezone,
      FirstPhotoDay,
      FirstPhotoMonth,
      FirstPhotoYear,
      FirstPhotoLocalTime,
      FirstPhotoAltitude,
      FirstPhotoLatitude,
      FirstPhotoLongitude,
      FirstPhotoCountry,
      FirstPhotoIso,
      FirstPhotoExposure,
      FirstPhotoFnumber,
      FirstPhotoFocalLength,
      FirstPhotoSubject,
      FirstPhotoArtist,
      FirstPhotoCopyright,
      FirstPhotoLicense,
      FirstPhotoDescription,
      FirstPhotoKeywords,
      FirstPhotoNotes
    );
    await contextmenu.checkContextMenuCount("1");
    await contextmenu.clearSelection();
  }
);

test.meta("testID", "photosgeo64-002").meta({ type: "short", mode: "public" })(
  "Common: Add 6 decimal Location",
  async (t) => {
    await menu.openPage("browse");
    await toolbar.setFilter("view", "Cards");
    const FirstPhotoUid = await photo.getNthPhotoUid("image", 0);
    await t.click(page.cardTitle.withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.latitude.visible).ok();

    await t.click(photoedit.dialogNext);

    await t.expect(photoedit.dialogPrevious.getAttribute("disabled")).notEql("disabled");

    await t.click(photoedit.dialogPrevious).click(photoedit.dialogClose);
    await photoviewer.openPhotoViewer("uid", FirstPhotoUid);
    await photoviewer.triggerPhotoViewerAction("edit");
    const FirstPhotoTitle = await photoedit.title.value;
    const FirstPhotoLocalTime = await photoedit.localTime.value;
    const FirstPhotoDay = await photoedit.day.value;
    const FirstPhotoMonth = await photoedit.month.value;
    const FirstPhotoYear = await photoedit.year.value;
    const FirstPhotoTimezone = await photoedit.timezone.value;
    const FirstPhotoLatitude = await photoedit.latitude.value;
    const FirstPhotoLongitude = await photoedit.longitude.value;
    const FirstPhotoAltitude = await photoedit.altitude.value;
    const FirstPhotoCountry = await photoedit.country.value;
    const FirstPhotoCamera = await photoedit.camera.innerText;
    const FirstPhotoIso = await photoedit.iso.value;
    const FirstPhotoExposure = await photoedit.exposure.value;
    const FirstPhotoLens = await photoedit.lens.innerText;
    const FirstPhotoFnumber = await photoedit.fnumber.value;
    const FirstPhotoFocalLength = await photoedit.focallength.value;
    const FirstPhotoSubject = await photoedit.subject.value;
    const FirstPhotoArtist = await photoedit.artist.value;
    const FirstPhotoCopyright = await photoedit.copyright.value;
    const FirstPhotoLicense = await photoedit.license.value;
    const FirstPhotoDescription = await photoedit.description.value;
    const FirstPhotoKeywords = await photoedit.keywords.value;
    const FirstPhotoNotes = await photoedit.notes.value;

    await t
      .typeText(photoedit.title, "Not saved photo title", { replace: true })
      .click(photoedit.detailsClose)
      .click(Selector("button.action-title-edit").withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.title.value).eql(FirstPhotoTitle);

    await photoedit.editPhoto(
      "New Photo Title",
      "Europe/Moscow",
      "15",
      "07",
      "2019",
      "05:30:30",
      "-1",
      "41.153334",
      "20.168331",
      "32",
      "1/32",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      ", cat, love",
      "Some notes"
    );
    if (t.browser.platform === "mobile") {
      await t.eval(() => location.reload());
    } else {
      await toolbar.triggerToolbarAction("reload");
    }
    await toolbar.search("uid:" + FirstPhotoUid);

    await t
      .expect(page.cardTitle.withAttribute("data-uid", FirstPhotoUid).innerText)
      .eql("New Photo Title");

    await photo.triggerHoverAction("uid", FirstPhotoUid, "select");
    await contextmenu.triggerContextMenuAction("edit", "");

    //const expectedValues = [{ FirstPhotoTitle: photoedit.title }, { "bluh bla": photoedit.day }];
    /*const expectedValues = [
    [FirstPhotoTitle, photoedit.title],
    ["blah", photoedit.day],
  ];
  await photoedit.checkEditFormValuesNewNew(expectedValues);*/

    await photoedit.checkEditFormValues(
      "New Photo Title",
      "15",
      "07",
      "2019",
      "05:30:30",
      "Europe/Moscow",
      "Albania",
      "-1",
      "41.153334",
      "20.168331",
      "",
      "32",
      "1/32",
      "",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      "cat",
      "Some notes"
    );

    await photoedit.undoPhotoEdit(
      FirstPhotoTitle,
      FirstPhotoTimezone,
      FirstPhotoDay,
      FirstPhotoMonth,
      FirstPhotoYear,
      FirstPhotoLocalTime,
      FirstPhotoAltitude,
      FirstPhotoLatitude,
      FirstPhotoLongitude,
      FirstPhotoCountry,
      FirstPhotoIso,
      FirstPhotoExposure,
      FirstPhotoFnumber,
      FirstPhotoFocalLength,
      FirstPhotoSubject,
      FirstPhotoArtist,
      FirstPhotoCopyright,
      FirstPhotoLicense,
      FirstPhotoDescription,
      FirstPhotoKeywords,
      FirstPhotoNotes
    );
    await contextmenu.checkContextMenuCount("1");
    await contextmenu.clearSelection();
  }
);

test.meta("testID", "photosgeo64-002").meta({ type: "short", mode: "public" })(
  "Common: Add 7 decimal Location",
  async (t) => {
    await menu.openPage("browse");
    await toolbar.setFilter("view", "Cards");
    const FirstPhotoUid = await photo.getNthPhotoUid("image", 0);
    await t.click(page.cardTitle.withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.latitude.visible).ok();

    await t.click(photoedit.dialogNext);

    await t.expect(photoedit.dialogPrevious.getAttribute("disabled")).notEql("disabled");

    await t.click(photoedit.dialogPrevious).click(photoedit.dialogClose);
    await photoviewer.openPhotoViewer("uid", FirstPhotoUid);
    await photoviewer.triggerPhotoViewerAction("edit");
    const FirstPhotoTitle = await photoedit.title.value;
    const FirstPhotoLocalTime = await photoedit.localTime.value;
    const FirstPhotoDay = await photoedit.day.value;
    const FirstPhotoMonth = await photoedit.month.value;
    const FirstPhotoYear = await photoedit.year.value;
    const FirstPhotoTimezone = await photoedit.timezone.value;
    const FirstPhotoLatitude = await photoedit.latitude.value;
    const FirstPhotoLongitude = await photoedit.longitude.value;
    const FirstPhotoAltitude = await photoedit.altitude.value;
    const FirstPhotoCountry = await photoedit.country.value;
    const FirstPhotoCamera = await photoedit.camera.innerText;
    const FirstPhotoIso = await photoedit.iso.value;
    const FirstPhotoExposure = await photoedit.exposure.value;
    const FirstPhotoLens = await photoedit.lens.innerText;
    const FirstPhotoFnumber = await photoedit.fnumber.value;
    const FirstPhotoFocalLength = await photoedit.focallength.value;
    const FirstPhotoSubject = await photoedit.subject.value;
    const FirstPhotoArtist = await photoedit.artist.value;
    const FirstPhotoCopyright = await photoedit.copyright.value;
    const FirstPhotoLicense = await photoedit.license.value;
    const FirstPhotoDescription = await photoedit.description.value;
    const FirstPhotoKeywords = await photoedit.keywords.value;
    const FirstPhotoNotes = await photoedit.notes.value;

    await t
      .typeText(photoedit.title, "Not saved photo title", { replace: true })
      .click(photoedit.detailsClose)
      .click(Selector("button.action-title-edit").withAttribute("data-uid", FirstPhotoUid));

    await t.expect(photoedit.title.value).eql(FirstPhotoTitle);

    await photoedit.editPhoto(
      "New Photo Title",
      "Europe/Moscow",
      "15",
      "07",
      "2019",
      "05:30:30",
      "-1",
      "41.1533348",
      "20.1683312",
      "32",
      "1/32",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      ", cat, love",
      "Some notes"
    );
    if (t.browser.platform === "mobile") {
      await t.eval(() => location.reload());
    } else {
      await toolbar.triggerToolbarAction("reload");
    }
    await toolbar.search("uid:" + FirstPhotoUid);

    await t
      .expect(page.cardTitle.withAttribute("data-uid", FirstPhotoUid).innerText)
      .eql("New Photo Title");

    await photo.triggerHoverAction("uid", FirstPhotoUid, "select");
    await contextmenu.triggerContextMenuAction("edit", "");

    //const expectedValues = [{ FirstPhotoTitle: photoedit.title }, { "bluh bla": photoedit.day }];
    /*const expectedValues = [
    [FirstPhotoTitle, photoedit.title],
    ["blah", photoedit.day],
  ];
  await photoedit.checkEditFormValuesNewNew(expectedValues);*/

    await photoedit.checkEditFormValues(
      "New Photo Title",
      "15",
      "07",
      "2019",
      "05:30:30",
      "Europe/Moscow",
      "Albania",
      "-1",
      "41.1533348",
      "20.1683312",
      "",
      "32",
      "1/32",
      "",
      "29",
      "33",
      "Super nice edited photo",
      "Happy",
      "Happy2020",
      "Super nice cat license",
      "Description of a nice image :)",
      "cat",
      "Some notes"
    );

    await photoedit.undoPhotoEdit(
      FirstPhotoTitle,
      FirstPhotoTimezone,
      FirstPhotoDay,
      FirstPhotoMonth,
      FirstPhotoYear,
      FirstPhotoLocalTime,
      FirstPhotoAltitude,
      FirstPhotoLatitude,
      FirstPhotoLongitude,
      FirstPhotoCountry,
      FirstPhotoIso,
      FirstPhotoExposure,
      FirstPhotoFnumber,
      FirstPhotoFocalLength,
      FirstPhotoSubject,
      FirstPhotoArtist,
      FirstPhotoCopyright,
      FirstPhotoLicense,
      FirstPhotoDescription,
      FirstPhotoKeywords,
      FirstPhotoNotes
    );
    await contextmenu.checkContextMenuCount("1");
    await contextmenu.clearSelection();
  }
);
