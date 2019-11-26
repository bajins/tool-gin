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
    arr.forEach(function (value) {
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
    arr.forEach(function (value) {
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
        if (obj.match(RegExp(`^.*${arr[i]}.*`))) {
            return true;
        }
    }
    return false;
}


/**
 * 类正态排序
 *
 * @param arr
 * @returns {[]}
 */
const normalSort = function (arr) {
    let temp = [];
    //先将数组从小到大排列得到 [1, 1, 2, 2, 3, 3, 3, 4, 6]
    let sortArr = arr.sort(function (a, b) {
        return a - b
    });
    for (let i = 0, l = arr.length; i < l; i++) {
        if (i % 2 == 0) {
            // 下标为偶数的顺序放到前边
            temp[i / 2] = sortArr[i];
        } else {
            // 下标为奇数的从后往前放
            temp[l - (i + 1) / 2] = sortArr[i];
        }
    }
    return temp;
}

/**
 * 利用Box-Muller方法极坐标形式
 * 使用两个均匀分布产生一个正态分布
 *
 * @param mean
 * @param sigma
 * @returns {*}
 */
const normalDistribution = function (mean, sigma) {
    let u = 0.0;
    let v = 0.0;
    let w = 0.0;
    let c;
    do {
        //获得两个（-1,1）的独立随机变量
        u = Math.random() * 2 - 1.0;
        v = Math.random() * 2 - 1.0;
        w = u * u + v * v;
    } while (w == 0.0 || w >= 1.0);

    c = Math.sqrt((-2 * Math.log(w)) / w);

    return mean + u * c * sigma;
}


/**
 * 随机拆分一个数
 *
 * @param total 总和
 * @param nums 个数
 * @param max 最大值
 * @returns {number[]}
 */
const randomSplit = function (total, nums, max) {
    let rest = total;
    let result = Array.apply(null, {length: nums}).map((n, i) => nums - i).map(n => {
        const v = 1 + Math.floor(Math.random() * (max | rest / n * 2 - 1));
        rest -= v;
        return v;
    });
    result[nums - 1] += rest;
    return result;
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
    trimSpace,
    trimFilter,
    notInArrayKV,
    inArrayKV,
    reinsertElement,
    mergeArray,
    splitArray,
    isInArray,
    normalSort,
    normalDistribution,
    randomSplit
}