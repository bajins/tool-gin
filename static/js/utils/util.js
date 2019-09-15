/**
 *
 * @Description:
 * @Author: claer
 * @File: util.js
 * @Version: 1.0.0
 * @Time: 2019/9/15 20:11
 * @Project: key-gin
 * @Package:
 * @Software: GoLand
 */


/**
 * 给String对象增加一个原型方法:
 * 判断一个字符串是以指定字符串结尾的
 *
 * @param str           需要判断的子字符串
 * @returns {boolean}   是否以该字符串结尾
 */
String.prototype.endWith = function (str) {
    if (str == null || str == "" || this.length == 0 || str.length > this.length)
        return false;
    if (this.substring(this.length - str.length) != str) {
        return false;
    }
    return true;
}


/**
 * 给String对象增加一个原型方法:
 * 判断一个字符串是以指定字符串开头的
 *
 * @param str           需要判断的子字符串
 * @returns {boolean}   是否以该字符串开头
 */
String.prototype.startWith = function (str) {
    if (str == null || str == "" || this.length == 0 || str.length > this.length)
        return false;
    if (this.substr(0, str.length) != str) {
        return false;
    }
    return true;
}

/**
 * 给String对象增加一个原型方法:
 * 判断一个字符串是以指定字符串结尾的
 *
 * @param str           需要判断的子字符串
 * @returns {boolean}   是否以该字符串结尾
 */
String.prototype.endWithRegExp = function (str) {
    let reg = new RegExp(str + "$");
    return reg.test(this);
}
/**
 * 给String对象增加一个原型方法:
 * 判断一个字符串是以指定字符串开头的
 *
 * @param str           需要判断的子字符串
 * @returns {boolean}   是否以该字符串开头
 */
String.prototype.startWithRegExp = function (str) {
    let reg = new RegExp("^" + str);
    return reg.test(this);
}


/**
 * 给String对象增加一个原型方法:
 * 替换全部字符串 - 无replaceAll的解决方案,自定义扩展js函数库
 * 原生js中并没有replaceAll方法，只有replace，如果要将字符串替换，一般使用replace
 *
 * @param FindText      要替换的字符串
 * @param RepText       新的字符串
 * @returns {string}
 */
String.prototype.replaceAll = function (FindText, RepText) {
    // g表示执行全局匹配，m表示执行多次匹配
    let regExp = new RegExp(FindText, "gm");
    return this.replace(regExp, RepText);
}

/**
 * 给Date对象增加一个原型方法：格式化
 *
 * @param fmt
 * @returns {void | string}
 */
Date.prototype.format = function (fmt) {
    let o = {
        "M+": this.getMonth() + 1,
        "d+": this.getDate(),
        "h+": this.getHours(),
        "m+": this.getMinutes(),
        "s+": this.getSeconds(),
        "q+": Math.floor((this.getMonth() + 3) / 3),
        "S": this.getMilliseconds()
    };
    if (/(y+)/.test(fmt)) {
        fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    }
    for (let k in o) {
        if (new RegExp("(" + k + ")").test(fmt)) {
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
        }
    }
    return fmt;
}

// JS获取当前文件所在的文件夹全路径
// let js = document.scripts;
// js = js[js.length - 1].src.substring(0, js[js.length - 1].src.lastIndexOf("/") + 1);

/**
 * 获取当前路径
 *
 * @returns {string}
 */
const getCurrAbsPath = () => {
    let a = {};
    let rExtractUri = /((?:http|https|file):\/\/.*?\/[^:]+)(?::\d+)?:\d+/;
    // let expose = +new Date();
    // let isLtIE8 = ('' + doc.querySelector).indexOf('[native code]') === -1;


    // FF,Chrome
    if (document.currentScript) {
        return document.currentScript.src;
    }

    let stack;
    try {
        a.b();
    } catch (e) {
        stack = e.fileName || e.sourceURL || e.stack || e.stacktrace;
    }
    // IE10
    if (stack) {
        let absPath = rExtractUri.exec(stack)[1];
        if (absPath) {
            return absPath;
        }
    }

    // IE5-9
    // for (let scripts = doc.scripts, i = scripts.length - 1, script; script = scripts[i--];) {
    //     if (script.className != expose && script.readyState === 'interactive') {
    //         script.className = expose;
    //         // if less than ie 8, must get abs path by getAttribute(src, 4)
    //         return isLtIE8 ? script.getAttribute('src', 4) : script.src;
    //     }
    // }
}

/**
 * 获取绝对路径
 *
 * @returns {string}
 */
const getPath = () => {
    let jsPath = document.currentScript ? document.currentScript.src : function () {
        let js = document.scripts,
            last = js.length - 1,
            src;
        for (let i = last; i > 0; i--) {
            if (js[i].readyState === 'interactive') {
                src = js[i].src;
                break;
            }
        }
        return src || js[last].src;
    }();
    return jsPath.substring(jsPath.lastIndexOf('/') + 1, jsPath.length);
}


/**
 * 生成从最小数到最大数的随机数
 * minNum 最小数
 * maxNum 最大数
 */
const randomNum = (minNum, maxNum) => {
    return parseInt(Math.random() * (maxNum - minNum + 1) + minNum, 10);
}


/**
 * 设置延时后再执行下一步操作
 *
 * @return
 * @Description
 * @author claer woytu.com
 * @date 2019/7/4 20:22
 */
const delay = ms => new Promise(resolve => setTimeout(resolve, ms));


/**
 * 判断js数组/对象是否为空
 * isPrototypeOf() 验证一个对象是否存在于另一个对象的原型链上。即判断 Object 是否存在于 $obj 的原型链上。
 * js中一切皆对象，也就是说，Object 也存在于数组的原型链上，因此这里数组需要先于对象检验。
 * Object.keys() 返回一个由给定对象的自身可枚举属性组成的数组，数组中属性名的排列顺序和使用 for...in 循环遍历该对象时返回的顺序一致
 *
 * @param $obj
 * @return {boolean}
 */
function isEmpty($obj) {
    // 找不到属性
    if (typeof ($obj) == 'undefined') {
        return true;
    }
    // 检验非数组/对象类型  EX：undefined   null  ''  根据自身要求添加其他适合的为空的值  如：0 ,'0','  '  等
    if ($obj === 0 || $obj === '' || $obj === null) {
        return true;
    }
    if (typeof ($obj) === "string") {
        // 移除字符串中所有 ''
        $obj = $obj.replace(/\s*/g, "");
        if ($obj === '') {
            return true;
        }
    }
    if (typeof ($obj) === "object") {
        if (!Array.isArray($obj) || $obj.length <= 0 || Object.keys($obj).length <= 0) {
            return true;
        }
    }
    return false;
}


