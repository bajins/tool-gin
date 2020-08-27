/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: detection.js
 * @Version: 1.0.0
 * @Time: 2019/11/26/026 15:49
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */

$(function () {
    let area_width = "500px";
    if (IEVersion() != -1) {
        const html = `<div style="font-weight:bold;text-align: center;padding: 20px;font-size: 200%">
                        不支持IE，请使用其他浏览器
                    </div>`;
        //自定页
        layer.open({
            // 在默认状态下，layer是宽高都自适应的，但当你只想定义宽度时，你可以area: '500px'，高度仍然是自适应的。
            // 当你宽高都要定义时，你可以area: ['500px', '300px']
            area: [area_width],
            type: 1,
            title: false,
            content: html,
            scrollbar: false,
            closeBtn: 0
        });
    }
    /**
     * 监听窗口变化
     */
    window.onresize = function () {
        if (window.innerWidth <= 600) {
            area_width = "80%";
        }
    }
})

/**
 * 判断IE以及Edge浏览器的版本
 *
 * @returns {string|number}
 * @constructor
 */
function IEVersion() {
    // 取得浏览器的userAgent字符串
    var userAgent = navigator.userAgent;
    // 判断是否IE<11浏览器
    if (userAgent.indexOf("compatible") > -1 && userAgent.indexOf("MSIE") > -1) {
        var reIE = new RegExp("MSIE (\\d+\\.\\d+);");
        reIE.test(userAgent);
        var fIEVersion = parseFloat(RegExp["$1"]);
        if (fIEVersion == 7) {
            return 7;
        } else if (fIEVersion == 8) {
            return 8;
        } else if (fIEVersion == 9) {
            return 9;
        } else if (fIEVersion == 10) {
            return 10;
        }
        // IE版本<=7
        else {
            return 6;
        }
    }
    // 判断是否IE的Edge浏览器
    else if (userAgent.indexOf("Edge") > -1 && ("ActiveXObject" in window)) {
        return 'edge';//edge
    }
    // IE11
    else if (userAgent.indexOf('Trident') > -1 && userAgent.indexOf("rv:11.0") > -1) {
        return 11;
    }
    // 不是ie浏览器
    else {
        return -1;
    }
}