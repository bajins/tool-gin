/**
 *
 * @Description:
 * @Author: claer
 * @File: string.js
 * @Version: 1.0.0
 * @Time: 2019/9/15 20:03
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */


/**
 * 生成一个指定长度的随机字符串
 *
 * @param len 指定长度
 * @param str 指定字符串范围，默认小写字母、数字、下划线
 * @returns {string}
 */
const randomString=(len, str)=> {
    str = str || 'abcdefghijklmnopqrstuvwxyz0123456789_';
    let randomString = '';
    for (let i = 0; i < len; i++) {
        let randomPoz = Math.floor(Math.random() * str.length);
        randomString += str.substring(randomPoz, randomPoz + 1);
    }
    return randomString;
}


/**
 * 正则表达式去除空行
 *
 * @param oldStr
 * @returns {string}
 */
const replaceBlank=(oldStr)=> {
    if (typeof oldStr != "string") {
        throw new Error("正则表达式去除空行，传入的不为字符串！");
    }
    // 匹配空行
    let reg = /\n(\n)*( )*(\n)*\n/g;
    return oldStr.replace(reg, "\n");
}


/**
 * 格式化数字为字符串
 *
 * @param n
 * @returns {string}
 */
const formatNumber = (n) => {
    n = n.toString();
    return n[1] ? n : '0' + n;
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
    randomString,
    replaceBlank,
    formatNumber
}