import { Selector, t } from "testcafe";

export default class Page {
  constructor() {
    this.placesSearch = Selector('input[placeholder="Search"]');
    this.openClusterInSearch = Selector("div.cluster-control-container button.action-browse");
    this.closeCluster = Selector("div.cluster-control-container button.action-close");
    this.clearLocation = Selector("button.action-clear-location");
    this.zoomOut = Selector("button.maplibregl-ctrl-zoom-out");
  }

  async search(term) {
    await t.typeText(this.placesSearch, term, { replace: true }).pressKey("enter");
  }
}
