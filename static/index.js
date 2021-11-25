(function () {
  /*----------- utilities type and function ----------*/
  const goGet = async (url) => {
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

  /*--------------- Main New Todo --------------*/

  const todoInput = document.getElementById("new-todo");

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

  const pendingTagDiv = document.getElementById("pending-tag");
  const tagInput = document.getElementById("input-tag");
  const newTagPendings = PendingTagsMap();
  let showSelectColor = false;

  const cancleTagEvent = (tagId) => (event) => {
    event.stopPropagation();
    removeChildElementByIndex(pendingTagDiv, newTagPendings.getIndexOf(tagId));
    newTagPendings.delete(tagId);
  }

  const selectTagColorEvent = (tagId, colorId) => (event) => {
    newTagPendings.selectColor(tagId, colorId);
    const indexElement = pendingTagDiv.children[newTagPendings.getIndexOf(tagId)];
    const hexColor = highlightColors.get(colorId).hex;
    indexElement.style.background = hexColor;
  }

  const createColorSelector = (idObj, name, hex) => {
    const color = document.createElement("div");
    const colorName = document.createElement("span");
    colorName.textContent = name;
    color.classList.add("color-selector");
    color.style.background = hex;
    color.addEventListener("click",
      selectTagColorEvent(idObj.tagId, idObj.colorId));
    return color;
  }

  const createTagHighlightSelectors = (tagId) => {
    const wrapper = document.createElement("div");
    const tagName = document.createElement("p");
    const colorWrapper = document.createElement("div");
    tagName.textContent = newTagPendings.get(tagId).name;
    const x = createCancleX();
    x.addEventListener("click", (event) => {
      showSelectColor = false;
      event.target.parentNode.remove();
    });
    wrapper.classList.add("highlight-selectors-wrapper");
    wrapper.append(tagName, colorWrapper);
    highlightColors.forEach((color, id) => {
      const ids = {tagId: tagId, colorId: id};
      colorWrapper.appendChild(createColorSelector(ids, color.name, color.hex));
    });
    return wrapper;
  }

  const showHighlightSelectorEvent = (tagId) => (event) => {
    if (showSelectColor) {
      return;
    }
    showSelectColor = true;
    event.target.after(createTagHighlightSelectors(tagId));
  }

  const createTagPendingElement = (id, name) => {
    const tagWrapper = document.createElement("div");
    tagWrapper.classList.add("tag-pending");
    const tagName = document.createElement("span");
    tagName.textContent = name;
    const cancleX = createCancleX();
    cancleX.addEventListener("click", cancleTagEvent(id));
    tagWrapper.appendChild(tagName);
    tagWrapper.appendChild(cancleX);
    tagWrapper.addEventListener("click", showHighlightSelectorEvent(tagId));
    return tagWrapper;
  }

  const suggestTags = (value) => {
    const suggestArray = [];
    availableTags.forEach((tag, id) => {
      if (tag.name.toLowerCase().includes(value.toLowerCase())) {
        suggestArray.push(id);
      }
    });
    return suggestArray;
  }

  const assignNewTagEvent = (tagId) => (event) => {
    const tag = availableTags.get(tagId);
    newTagPendings.setId(tagId, {name: tag.name, color: 0});
    pendingTagDiv.appendChild(createTagPendingElement(tagId, tag.name));
    event.target.parentNode.remove();
  }

  const createTagSuggestListItem = (tagId) => {
    const li = document.createElement("li");
    li.textContent = availableTags.get(tagId).name;
    li.addEventListener("click", assignNewTagEvent(tagId));
    return li;
  }

  const createTagSuggestList = (tags) => {
    const ul = document.createElement("ul");
    ul.classList.add("tag-suggest-list");
    tags.forEach((tagId) => ul.appendChild(createTagSuggestListItem(tagId)));
    return ul;
  }

  const tagSearchEvent = (event) => {
    const tags = suggestTags(event.target.value);
  }

  tagInput.addEventListener("focus", (event) => {

  });

  tagInput.addEventListener("focusout", (event) => {

  })

  tagInput.addEventListener("change", tagSearchEvent);

})();
