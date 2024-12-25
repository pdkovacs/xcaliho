import type { ExcalidrawElement } from "@excalidraw/excalidraw/types/element/types";

export interface Drawing {
  readonly title: string;
  readonly content: string;
}

const mockDrawingList = [
  "kalap kabat",
  "gubanc bucka"
];

const mockDrawing: ExcalidrawElement = {
  "id": "zTROlT6-QlvQKxH9byf1m",
  "type": "rectangle",
  "x": 552,
  "y": 185.625,
  "width": 160,
  "height": 77.5,
  "angle": 0,
  "strokeColor": "#1e1e1e",
  "backgroundColor": "transparent",
  "fillStyle": "solid",
  "strokeWidth": 2,
  "strokeStyle": "solid",
  "roughness": 1,
  "opacity": 100,
  "groupIds": [],
  "frameId": null,
  "roundness": {
      "type": 3
  },
  "seed": 462881513,
  "version": 13,
  "versionNonce": 305270121,
  "isDeleted": false,
  "boundElements": null,
  "updated": 1735254351553,
  "link": null,
  "locked": false
};

export const fetchDrawingList = () => {
  return new Promise<{drawingList: string[]}>(resolve =>
    setTimeout(() => resolve({ drawingList: mockDrawingList }), 1000)
  );
};

export const fetchDrawing = (title: string) => {
  console.info("Fetching drawing: ", title);
  return new Promise<Drawing>(resolve =>
    setTimeout(() => resolve({ title, content: JSON.stringify([mockDrawing]) }), 1000)
  );
};

export const saveDrawing = (name: string, content: string) => {
  console.info("Saving drawing of length ", content.length, ", as ", name);
  return new Promise<void>(resolve =>
    setTimeout(() => resolve(), 1000)
  );
};
