import BaseEvent from '../BaseEvent';
import XType from '../../ui/XType';

/**
 * 历史面板改变事件
 * @param {*} app 
 */
function HistoryPanelEvent2(app) {
    BaseEvent.call(this, app);

    this.currentObject = null;
    this.copiedMaterial = null;
}

HistoryPanelEvent2.prototype = Object.create(BaseEvent.prototype);
HistoryPanelEvent2.prototype.constructor = HistoryPanelEvent2;

HistoryPanelEvent2.prototype.start = function () {
    this.app.on(`editorCleared.${this.id}`, this.refreshUI.bind(this));
    this.app.on(`historyChanged.${this.id}`, this.refreshUI.bind(this));
    this.refreshUI();
};

HistoryPanelEvent2.prototype.stop = function () {
    this.app.on(`editorCleared.${this.id}`, null);
    this.app.on(`historyChanged.${this.id}`, null);
};

var refreshUI = function () {
    var editor = this.app.editor;
    var history = editor.history;
    var outliner = XType.getControl('historyOutlinear');

    var options = [];
    var enumerator = 1;

    function buildOption(object) {
        var option = document.createElement('div');
        option.value = object.id;
        return option;
    }

    (function addObjects(objects) {
        for (var i = 0, l = objects.length; i < l; i++) {
            var object = objects[i];
            var option = buildOption(object);
            option.innerHTML = '&nbsp;' + object.name;
            options.push(option);
        }
    })(history.undos);


    (function addObjects(objects, pad) {
        for (var i = objects.length - 1; i >= 0; i--) {
            var object = objects[i];
            var option = buildOption(object);
            option.innerHTML = '&nbsp;' + object.name;
            option.style.opacity = 0.3;
            options.push(option);
        }
    })(history.redos, '&nbsp;');

    outliner.setOptions(options);
};


export default HistoryPanelEvent2;