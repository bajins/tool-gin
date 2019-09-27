/**
 *
 * @Description:
 * @Author: claer
 * @File: array.js
 * @Version: 1.0.0
 * @Time: 2019/9/15 20:01
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */


/**
 * splice方法删除数组中的空值
 *
 * @param array
 * @returns {*}
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
 */
function trimFilter(array) {
    array.filter(function (s) {
        return s && s.trim(); // 注：IE9(不包含IE9)以下的版本没有trim()方法
    });
}


/**
 * 过滤不在数组中的值
 *
 * @param arr           元数据数组
 * @param retentionArr  需要保留的值数组
 * @returns {[]}        去掉值后的新数组
 */
function notInArrayKV(arr, retentionArr) {
    let newArr = [];
    arr.forEach(function (value, index, array) {
        // 判断文件名以什么开头、是否在指定数组中存在
        if (!value.startsWith(".") && !retentionArr.includes(value)) {
            newArr.push(value);
        }
    });
    return newArr;
}


/**
 * 过滤在数组中的值
 *
 * @param arr           元数据数组
 * @param ignoresArr    需要去除的值数组
 * @returns {[]}        去掉值后的新数组
 */
function inArrayKV(arr, ignoresArr) {
    let newArr = [];
    arr.forEach(function (value, index, array) {
        // 判断文件名以什么开头、是否在指定数组中存在
        if (!value.startsWith(".") && ignoresArr.includes(value)) {
            newArr.push(value);
        }
    });
    return newArr;
}

/**
 * 插入去重的元素
 *
 * @param array
 * @param element
 * @returns {*}
 */
const reinsertElement = (array, element) => {
    if (array.indexOf(element) === -1) {
        array.push(element);
    }
    return array;
}


/**
 * 自定义数组合并并去重函数
 *
 * @param arr1
 * @param arr2
 * @returns {*}
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
 * 将数组平均分割
 *
 * @param arr   数组
 * @param len   分割成多少个
 * @returns {[]}
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
 * 判断数组中是否包含指定字符串
 *
 * @param arr
 * @param obj
 * @returns {boolean}
 */
function isInArray(arr, obj) {
    let i = arr.length;
    while (i--) {
        if (obj.match(RegExp("^.*" + arr[i] + ".*"))) {
            return true;
        }
    }
    return false;
}