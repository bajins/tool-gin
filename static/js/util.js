/**
 groovy/lang/GroovyShell
 * @return
 * @Description 获取当前路径
 * @author claer woytu.com
 * @date 2019/4/29 13:28
 */
const getCurrAbsPath = () => {
    let a = {},
        rExtractUri = /((?:http|https|file):\/\/.*?\/[^:]+)(?::\d+)?:\d+/;
    // expose = +new Date(),
    // isLtIE8 = ('' + doc.querySelector).indexOf('[native code]') === -1;


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
 groovy/lang/GroovyShell
 * @return
 * @Description 获取绝对路径
 * @author claer woytu.com
 * @date 2019/4/29 13:31
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
 * 时间转换工具
 * date 时间
 * join 年月日之间连接的字符
 */
const formatTime = (date, join) => {
    let year = date.getFullYear();
    let month = date.getMonth() + 1;
    let day = date.getDate();
    let hour = date.getHours();
    let minute = date.getMinutes();
    let second = date.getSeconds();
    return [year, month, day].map(formatNumber).join(join) + ' ' + [hour, minute, second].map(formatNumber).join(':');
}

/**
 *
 */
const formatNumber = (n) => {
    n = n.toString();
    return n[1] ? n : '0' + n;
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
 * 将数组平均分割
 * arr 数组
 * len 分割成多少个
 */
const splitArray = (arr, len) => {
    let arr_length = arr.length;
    let newArr = [];
    for (let i = 0; i < arr_length; i += len) {
        newArr.push(arr.slice(i, i + len));
    }
    return newArr;
}

/**
 * 自定义数组合并并去重函数
 * @param arr1 数组1
 * @param arr2 数组2
 * @return
 * @Description 自定义数组合并并去重函数
 * @author claer woytu.com
 * @date 2019/4/29 20:10
 */
const mergeArray = (arr1, arr2) => {
    // var _arr = new Array();
    // for (var i = 0; i < arr1.length; i++) {
    //   _arr.push(arr1[i]);
    // }
    // for (var i = 0; i < arr2.length; i++) {
    //   var flag = true;
    //   for (var j = 0; j < arr1.length; j++) {
    //     if (arr2[i] == arr1[j]) {
    //       flag = false;
    //       break;
    //     }
    //   }
    //   if (flag) {
    //     _arr.push(arr2[i]);
    //   }
    // }

    for (let i = 0; i < arr2.length; i++) {
        if (arr1.indexOf(arr2[i]) === -1) {
            arr1.push(arr2[i]);
        }
    }
    return arr1;
}

/**
 * 插入去重的元素
 *
 * @return
 * @Description
 * @author claer woytu.com
 * @date 2019/4/30 17:05
 */
const reinsertElement = (array, element) => {
    if (array.indexOf(element) === -1) {
        array.push(element);
    }
    return array;
}


/**
 * 判断js数组/对象是否为空
 * isPrototypeOf() 验证一个对象是否存在于另一个对象的原型链上。即判断 Object 是否存在于 $obj 的原型链上。
 * js中一切皆对象，也就是说，Object 也存在于数组的原型链上，因此这里数组需要先于对象检验。
 * Object.keys() 返回一个由给定对象的自身可枚举属性组成的数组，数组中属性名的排列顺序和使用 for...in 循环遍历该对象时返回的顺序一致
 * @param $obj
 * @return {boolean}
 * @Description
 * @author claer woytu.com
 * @date 2019/4/29 20:12
 */
const isEmpty = ($obj) => {
    // 找不到属性
    if (typeof $obj == 'undefined') {
        return true;
    }
    // 检验非数组/对象类型  EX：undefined   null  ''  根据自身要求添加其他适合的为空的值  如：0 ,'0','  '  等
    if ($obj === 0 || $obj === '' || $obj === null) {
        return true;
    }
    if (typeof $obj === "string") {
        $obj = $obj.trim().replace(/\s*/g, ""); //移除字符串中所有 ''
        if ($obj === '') {
            return true;
        }
    } else if (typeof $obj === "object") {
        if (!Array.isArray($obj) || $obj.length <= 0) {
            return true;
        }
        if (!Object.prototype.isPrototypeOf($obj) || !Object.keys($obj).length != 0) {
            return true;
        }
    }
    return false;
}

/**
 * replace默认只替换第一个匹配项
 * @param str 父字符串
 * @param substring 被替换的字符串
 * @param newString 新字符串
 *
 * "g"是匹配全部的意思，也可以换成""，就是匹配第一个
 *
 * @return
 * @Description
 * @author claer woytu.com
 * @date 2019/4/30 15:20
 */
const replace = (str, substring, newString, isAll) => {
    if (!isEmpty(isAll) && isAll) {
        return str.replace(new RegExp(substring, "g"), newString);
    }
    return str.replace(new RegExp(substring, ""), newString);
}

/**
 * 正则表达式去除空行
 *
 * @param oldStr 字符串
 * @return
 * @Description
 * @author claer woytu.com
 * @date 2019/6/13 17:55
 */
function replaceBlank(oldStr) {
    if (typeof oldStr != "string") {
        console.log("正则表达式去除空行，传入的不为字符串！");
    } else {
        // 匹配空行
        let reg = /\n(\n)*( )*(\n)*\n/g;
        return oldStr.replace(reg, "\n");
    }
}

/**
 * splice方法删除数组中的空值
 *
 * @param array
 * @return
 * @Description
 * @author claer woytu.com
 * @date 2019/6/13 18:14
 */
function trimSpace(array) {
    for (let i = 0; i < array.length; i++) {
        if (array[i] == " " || array[i] == null || typeof (array[i]) == "undefined") {
            array.splice(i, 1);
            i = i - 1;
        }
    }
    return array;
}

/**
 * filter 过滤方法删除数组中的空值
 *
 * @param array
 * @return
 * @Description
 * @author claer woytu.com
 * @date 2019/6/13 18:14
 */
function trimFilter(array) {
    array.filter(function (s) {
        return s && s.trim(); // 注：IE9(不包含IE9)以下的版本没有trim()方法
    });
}