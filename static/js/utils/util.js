/**
 *
 * @Description:
 * @Author: claer
 * @File: util.js
 * @Version: 1.0.0
 * @Time: 2019/9/15 20:11
 * @Project: tool-gin
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


if (!String.prototype.trim) {
    String.prototype.trim = function () {
        return this.replace(/^[\s\uFEFF\xA0]+|[\s\uFEFF\xA0]+$/g, '');
    }
}

if (!String.prototype.startsWith) {
    String.prototype.startsWith = function (searchString, position) {
        position = position || 0;
        return this.substr(position, searchString.length) === searchString;
    }
}


if (!String.prototype.endsWith) {
    String.prototype.endsWith = function (searchString, position) {
        let subjectString = this.toString();
        if (typeof position !== 'number' || !isFinite(position) || Math.floor(position) !== position || position > subjectString.length) {
            position = subjectString.length;
        }
        position -= searchString.length;
        let lastIndex = subjectString.indexOf(searchString, position);
        return lastIndex !== -1 && lastIndex === position;
    }
}


if (!String.prototype.includes) {
    String.prototype.includes = function (search, start) {
        'use strict';
        if (typeof start !== 'number') {
            start = 0;
        }

        if (start + search.length > this.length) {
            return false;
        } else {
            return this.indexOf(search, start) !== -1;
        }
    }
}

if (!String.prototype.repeat) {
    String.prototype.repeat = function (count) {
        if (this == null) {
            throw new TypeError('can\'t convert ' + this + ' to object');
        }
        let str = '' + this;
        count = +count;
        if (count != count) {
            count = 0;
        }
        if (count < 0) {
            throw new RangeError('repeat count must be non-negative');
        }
        if (count == Infinity) {
            throw new RangeError('repeat count must be less than infinity');
        }
        count = Math.floor(count);
        if (str.length == 0 || count == 0) {
            return '';
        }
        // Ensuring count is a 31-bit integer allows us to heavily optimize the
        // main part. But anyway, most current (August 2014) browsers can't handle
        // strings 1 << 28 chars or longer, so:
        if (str.length * count >= 1 << 28) {
            throw new RangeError('repeat count must not overflow maximum string size');
        }
        let rpt = '';
        for (; ;) {
            if ((count & 1) == 1) {
                rpt += str;
            }
            count >>>= 1;
            if (count == 0) {
                break;
            }
            str += str;
        }
        // Could we try:
        // return Array(count + 1).join(this);
        return rpt;
    }
}

//removes element from array
if (!Array.prototype.remove) {
    Array.prototype.remove = function (index) {
        this.splice(index, 1);
    }
}


if (!String.prototype.contains) {
    String.prototype.contains = String.prototype.includes;
}

if (!Array.prototype.insert) {
    Array.prototype.insert = function (index, item) {
        this.splice(index, 0, item);
    }
}


// ======================================  全局工具函数  =======================================


// JS获取当前文件所在的文件夹全路径
// let js = document.scripts;
// js = js[js.length - 1].src.substring(0, js[js.length - 1].src.lastIndexOf("/") + 1);

/**
 * 获取当前路径
 *
 * @returns {string}
 */
const getCurrAbsPath = () => {

    // ECMAScript6
    if (import.meta) {
        return import.meta.url;
    }

    // FF,Chrome
    if (document.currentScript) {
        return document.currentScript.src;
    }

    // IE10
    var e = new Error('');
    var stack = e.stack || e.sourceURL || e.stacktrace || '';
    var rExtractUri = /((?:http|https|file):\/\/.*?\/[^:]+)(?::\d+)?:\d+/;
    // var rgx = /(?:http|https|file):\/\/.*?\/.+?.js/;
    var absPath = rExtractUri.exec(stack);
    if (absPath) {
        return absPath[1];
    }

    // IE5-9
    var doc = exports.document;
    var scripts = doc.scripts;
    var expose = +new Date();
    var isLtIE8 = ('' + doc.querySelector).indexOf('[native code]') === -1;
    for (var i = 0; i < scripts.length; i++) {
        var script = scripts[i];
        if (script.className != expose && script.readyState === 'interactive') {
            script.className = expose;
            // 如果小于ie 8，则必须通过getAttribute（src，4）获得abs路径
            return isLtIE8 ? script.getAttribute('src', 4) : script.src;
        }
    }
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
 * 判断Array/Object/String是否为空
 *
 * @param obj
 * @return {boolean}
 */
const isEmpty = (obj) => {
    let type = Object.prototype.toString.call(obj);
    if (obj == null || obj == undefined) {
        return true;
    }
    switch (type) {
        case "[object Undefined]", "[object Null]":
            return true;
        case "[object String]":
            obj = obj.replace(/\s*/g, "");
            if (obj === "" || obj.length == 0) {
                return true;
            }
            return false;
        case "[object Array]":
            if (!Array.isArray(obj) || obj.length == 0) {
                return true;
            }
            return false;
        case "[object Object]":
            // Object.keys() 返回一个由给定对象的自身可枚举属性组成的数组
            if (obj.length == 0 || Object.keys(obj).length == 0) {
                return true;
            }
            return false;
        default:
            throw TypeError("只能判断Array/Object/String，当前类型为:" + type);
    }
}


/**
 * export default 服从 ES6 的规范,补充：default 其实是别名
 * module.exports 服从CommonJS 规范
 * 一般导出一个属性或者对象用 export default
 * 一般导出模块或者说文件使用 module.exports
 *
 * import from 服从ES6规范,在编译器生效
 * require 服从ES5 规范，在运行期生效
 * 目前 vue 编译都是依赖label 插件，最终都转化为ES5
 *
 * @return 将方法、变量暴露出去
 * @Description
 * @author claer woytu.com
 * @date 2019/4/29 11:58
 */
export default {
    getCurrAbsPath,
    getPath,
    randomNum,
    delay,
    isEmpty
}