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
    Object.keys(obj).forEach((key) => {
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
          console.error(`cannot set property. missing id field.`);
          return;
        }
        const entry = excludeId(current);
        if (!(this.validateValue(entry))) {
          console.error(`cannot set property. entry's value not valid.`);
          return;
        }
        this.set(current.id, entry);
      }
    }

    setId (id, value) {
      if (this.validateValue(value)) {
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

  const removeAllChildren = (parent) => {
    while (parent.firstChild) {
      parent.removeChild(parent.firstChild);
    }
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
      super();
      this.requireFields = ["name", "hex"];
    }
  }

  const availableTags = new ServerTagsMap();
  const highlightColors = new HighlightColorsMap();

  document.addEventListener("DOMContentLoaded", () => {
    goGet("/tags").then((tags) => availableTags.setAll(tags)).catch((e) => console.error(e));
    goGet("/highlight").then((colors) => highlightColors.setAll(colors)).catch((e) => console.error(e));
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

  const pendingTagDiv = document.getElementById("pending-tag");
  const tagInput = document.getElementById("input-tag");
  const newTagPendings = new PendingTagsMap();

  const getColorSelectorBox = () => {
    return document.getElementById(highlightSelectorBoxId);
  }

  const colorSelectorBoxIsOpen = () => {
    return getColorSelectorBox() !== null;
  }

  const closeColorSelectorBox = () => {
    if (!colorSelectorBoxIsOpen()) {
      return;
    }
    getColorSelectorBox().remove();
  }

  const cancleTagEvent = (tagId) => (event) => {
    event.stopPropagation();
    if (colorSelectorBoxIsOpen()) {
      const colorBoxCurrentTag = getColorSelectorBox().firstChild.textContent;
      const deleteTagName = newTagPendings.get(tagId).name;
      console.log(deleteTagName);
      console.log(colorBoxCurrentTag);
      if (colorBoxCurrentTag === deleteTagName) {
        closeColorSelectorBox();
      }
    }
    removeChildElementByIndex(pendingTagDiv, newTagPendings.getIndexOf(tagId));
    newTagPendings.delete(tagId);
  }

  const highlightSelectorBoxId = "color-selector-id";

  const selectTagColorEvent = (tagId, colorId) => (event) => {
    event.stopImmediatePropagation();
    newTagPendings.selectColor(tagId, colorId);
    const indexElement = pendingTagDiv.children[newTagPendings.getIndexOf(tagId)];
    const hexColor = highlightColors.get(colorId).hex;
    indexElement.style.background = hexColor;
  }

  const createColorSelector = (idObj, name, hex) => {
    const color = document.createElement("div");
    const colorName = document.createElement("span");
    colorName.textContent = name;
    color.style.width = "20px";
    color.style.height = "20px";
    color.classList.add("color-selector");
    color.style.background = hex;
    color.addEventListener("click",
      selectTagColorEvent(idObj.tagId, idObj.colorId));
    return color;
  }

  const createTagHighlightSelectors = (tagId) => {
    if (colorSelectorBoxIsOpen()) {
      throw new Error("can only have one color selector.");
    }
    const wrapper = document.createElement("div");
    wrapper.id = highlightSelectorBoxId;
    const tagName = document.createElement("p");
    const colorWrapper = document.createElement("div");
    tagName.textContent = newTagPendings.get(tagId).name;
    const x = createCancleX();
    x.addEventListener("click", (_) => closeColorSelectorBox());
    wrapper.classList.add("highlight-selectors-wrapper");
    wrapper.append(tagName, colorWrapper, x);
    highlightColors.forEach((color, id) => {
      const ids = {tagId: tagId, colorId: id};
      colorWrapper.appendChild(createColorSelector(ids, color.name, color.hex));
    });
    return wrapper;
  }

  const showHighlightSelectorEvent = (tagId) => (_) => {
    if (colorSelectorBoxIsOpen()) {
      return;
    }
    pendingTagDiv.after(createTagHighlightSelectors(tagId));
  }

  const createTagPendingElement = (tagId, name) => {
    const tagWrapper = document.createElement("div");
    tagWrapper.classList.add("tag-pending");
    const tagName = document.createElement("span");
    tagName.textContent = name;
    const cancleX = createCancleX();
    cancleX.addEventListener("click", cancleTagEvent(tagId));
    tagWrapper.appendChild(tagName);
    tagWrapper.appendChild(cancleX);
    tagWrapper.addEventListener("click", showHighlightSelectorEvent(tagId));
    return tagWrapper;
  }

  const suggestTags = (value, options) => {
    const suggestArray = [];
    availableTags.forEach((tag, id) => {
      if (options.excludeExistedTag && newTagPendings.get(id) !== undefined) {
        return;
      }
      if (tag.name.toLowerCase().includes(value.toLowerCase())) {
        suggestArray.push(id);
      }
    });
    return suggestArray;
  }

  const tagSuggestBoxId = "suggest-box";

  const getTagSuggestBox = () => {
    return document.getElementById(tagSuggestBoxId);
  }
  const tagSuggestBoxIsOpen = () => {
    return getTagSuggestBox() !== null;
  }

  const closeTagSuggest = () => {
    const searchBox = getTagSuggestBox();
    searchBox.remove();
  }

  const assignNewTagEvent = (tagId) => (event) => {
    event.stopImmediatePropagation();
    if (!newTagPendings.get(tagId)) {
      const tag = availableTags.get(tagId);
      newTagPendings.setId(tagId, {name: tag.name, color: 0});
      pendingTagDiv.appendChild(createTagPendingElement(tagId, tag.name));
      closeTagSuggest();
    }
  }

  const createTagSuggestListItem = (tagId) => {
    const li = document.createElement("li");
    li.textContent = availableTags.get(tagId).name;
    li.addEventListener("click", assignNewTagEvent(tagId));
    return li;
  }

  const appendTagsToList = (targetList, tagIds) => {
    tagIds.forEach((tagId) => targetList.appendChild(createTagSuggestListItem(tagId)));
  }

  const createTagSuggestList = (tags) => {
    if (tagSuggestBoxIsOpen()) {
      const error = new Error("can only have one tag suggest list.");
      console.error(error);
      return;
    }
    const ul = document.createElement("ul");
    ul.id = tagSuggestBoxId;
    ul.classList.add("tag-suggest-list");
    appendTagsToList(ul, tags);
    return ul;
  }

  const tagSearchInputEvent = (event) => {
    const input = event.target;
    let searchList = getTagSuggestBox();
    const tags = suggestTags(input.value, {excludeExistedTag: true});
    if (!getTagSuggestBox()) {
      searchList = createTagSuggestList(tags);
      input.after(searchList);
    } else {
      removeAllChildren(searchList);
      appendTagsToList(searchList, tags);
    }
  }

  const createTagSearchEvent = (event) => {
    if (tagSuggestBoxIsOpen()) {
      return;
    }
    const input = event.target;
    const allTag = suggestTags("", {excludeExistedTag: true});
    const suggestList = createTagSuggestList(allTag);
    input.after(suggestList);
  }
  
  const closeSearchBoxEvent = (_) => {
    if (tagSuggestBoxIsOpen()) {
      closeTagSuggest();
    }
  }

  tagInput.addEventListener("click", (event) => event.stopImmediatePropagation());
  tagInput.addEventListener("focus", createTagSearchEvent);
  tagInput.addEventListener("input", tagSearchInputEvent);
  document.addEventListener("click", closeSearchBoxEvent);

  /*--------------- Main New Todo --------------*/

  const addForm = document.getElementById("add-form");
  const tagInputWrapper = document.getElementById("tag-inputW");
  const todoInput = document.querySelector("input[name=new-todo]");
  const todoSubmitBtn = document.getElementById("todoBtn");

  tagInputWrapper.addEventListener("submit", (event) => event.preventDefault());

  const disableForm = (disabled=true) => {
    tagInput.disabled = disabled;
    todoInput.disabled = disabled;
    todoSubmitBtn.disabled = disabled;
    pendingTagDiv.style.pointerEvents = disabled ? "none" : "auto";
  }

  const submitTodoEvent = (event) => {
    event.preventDefault();
    const data = new FormData(addForm);
    const body = {
      todo: data.get("new-todo"),
      tags: [],
    }
    newTagPendings.forEach((value, id) => {
      body.tags.push({
        tagId: id,
        name: value.name,
        colorId: value.color,
      });
    });
    disableForm();
    fetch("/", {
      method: "POST",
      headers: {"Content-type": "application/json"},
      body: JSON.stringify(body),
    }).then((res) => res.json())
      .then((data) => {
        console.log(data);
      })
      .catch((err) => console.error(err))
      .finally(() => disableForm(false))
  }

  addForm.addEventListener("submit", submitTodoEvent);

})();
