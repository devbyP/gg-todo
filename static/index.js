(function () {
  /*----------- utilities type and function ----------*/
  const goGet = (url) => {
    return fetch(url).then((res) => res.json());
  }

  const hasNoId = (obj) => {
    return !obj.hasOwnProperty("id");
  }

  const excludeId = (obj) => {
    const nonId = {};
    obj.keys().forEach((key) => {
      if (key !== "id") {
        nonId[key] = obj[key];
      }
    })
    return nonId;
  }

  class IdMap extends Map {
    constructor () {
      super();
      this.requireFields = [];
    }

    validateValue (value) {
      let own = true;
      for (let i = 0; i < this.requireFields.length; i++) {
        own = value.hasOwnProperty(this.requireFields[i]);
      }
      return own;
    }

    setAll (idArr) {
      for (let i = 0; i < idArr.length; i++) {
        const current = idArr[i];
        if (hasNoId(current)) {
          console.error(`cannot set to ${this.name}. missing id field.`);
          return;
        }
        const entry = excludeId(current);
        if (!this.validateValue(entry)) {
          console.error(`cannot set to ${this.name}. entry's value not valid.`);
          return;
        }
        this.set(current.id, entry);
      }
    }

    setId (id, value) {
      if (validateValue(value)) {
        this.set(id, value);
      }
    }

    getIndexOf (id) {
      const keysIterator = this.keys();
      let currentIndex = 0;
      for (let result = keysIterator.next(); !result.done; currentIndex++ ) {
        if (result.value === id) {
          return currentIndex;
        }
        result = keysIterator.next();
      }
    }
  }

  /*------------- General Components ------------*/
  const createCancleX = () => {
    const cancleX = document.createElement("span");
    cancleX.textContent = "X";
    return cancleX;
  }

  /*--------------- DOM Functions -------------- */
  const removeChildElementByIndex = (wrapper, index) => {
    wrapper.children[index].remove();
  }

  /*--------------- Page Init ---------------*/
  class ServerTagsMap extends IdMap {
    constructor() {
      super();
      this.requireFields = ["name"];
    }
  }

  class HighlightColorsMap extends IdMap {
    constructor() {
      this.requireFields = ["name", "hex"];
    }
  }

  const availableTags = ServerTagsMap();
  const highlightColors = HighlightColorsMap();

  document.addEventListener("DOMContentLoaded", () => {
    goGet("/tags").then(availableTags.setAll).catch(console.error);
    goGet("/highlight").then(highlightColors.setAll).catch(console.error);
  });

  /*------------ Pending Tags --------------*/

  class PendingTagsMap extends IdMap {
    constructor () {
      super();
      this.requireFields = ["name", "color"];
    }
    
    selectColor (tagId, colorId) {
      const tag = this.get(tagId);
      tag.color = colorId;
      this.set(tagId, tag);
    }
  } 

  const tagSubmitButton = document.getElementById("add-tag");
  const pendingTagDiv = document.getElementById("pending-tag");
  const tagInput = document.getElementById("input-tag");
  const newTagPendings = PendingTagsMap();

  const cancleTagEvent = (tagId) => (event) => {
    event.stopPropagation();
    removeChildElementByIndex(pendingTagDiv, newTagPendings.getIndexOf(tagId));
    newTagPendings.delete(tagId);
  }

  const selectTagColorEvent = (tagId, colorId) => (event) => {
    newTagPendings.selectColor(tagId, colorId);
    const index = pendingTagDiv.children[newTagPendings.getIndexOf(tagId)];
    const hexColor = highlightColors.get(colorId).hex;
    index.style.background = hexColor;
  }

  const createColorSelector = (idObj, name, hex) => {
    const color = document.createElement("div");
    const colorName = document.createElement("span");
    colorName.textContent = name;
    color.classList.add("color-selector");
    color.style.background = hex;
    color.onclick(selectTagColorEvent(idObj.tagId, idObj.colorId));
    return color;
  }

  const createTagHighlightSelectors = (tagId) => {
    const wrapper = document.createElement("div");
    const tagName = document.createElement("p");
    const colorWrapper = document.createElement("div");
    tagName.textContent = newTagPendings.get(tagId).name;
    wrapper.classList.add("highlight-selectors-wrapper");
    wrapper.append(tagName, colorWrapper);
    highlightColors.forEach((color, id) => {
      const ids = {tagId: tagId, colorId: id};
      colorWrapper.appendChild(createColorSelector(ids, color.name, color.hex));
    })
    return wrapper;
  }

  const showHighlightSelectorEvent = (tagId) => (event) => {
    
  }

  const createTagPendingElement = (id, name) => {
    const tagWrapper = document.createElement("div");
    tagWrapper.classList.add("tag-pending");
    const tagName = document.createElement("span");
    tagName.textContent = name;
    const cancleX = createCancleX();
    cancleX.onclick(cancleTagEvent(id))
    tagWrapper.appendChild(tagName);
    tagWrapper.appendChild(cancleX);
    return tagWrapper;
  }


  const assignNewTag = (event) => {
    event.preventDefault();
  }

  tagSubmitButton.addEventListener("click", assignNewTag);

  /*--------------- Main New Todo --------------*/

  const todoInput = document.getElementById("new-todo");

})()
