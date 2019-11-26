/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: color.js
 * @Version: 1.0.0
 * @Time: 2019/11/21 21:14
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */

// 首字母配置颜色
const colorArr = {
    'a': 'rgb(17, 1, 65)',
    'b': 'rgb(113, 1, 98)',
    'c': 'rgb(161, 42, 94)',
    'd': 'rgb(237, 3, 69)',
    'e': 'rgb(239, 106, 50)',
    'f': 'rgb(251, 191, 69)',
    'g': 'rgb(170, 217, 98)',
    'h': 'rgb(3, 195, 131)',
    'i': 'rgb(1, 115, 81)',
    'j': 'rgb(1, 84, 90)',
    'k': 'rgb(38, 41, 74)',
    'l': 'rgb(26, 19, 52)',
    'm': 'rgb(0, 102, 119)',
    'n': 'rgb(119, 153, 85)',
    'o': 'rgb(255, 170, 102)',
    'p': 'rgb(255, 119, 119)',
    'q': 'rgb(199, 96, 101)',
    'r': 'rgb(23, 103, 87)',
    's': 'rgb(188, 173, 148)',
    't': 'rgb(83, 109, 114)',
    'u': 'rgb(102, 188, 41)',
    'v': 'rgb(181, 231, 146)',
    'w': 'rgb(232, 247, 221)',
    'x': 'rgb(113, 39, 122)',
    'y': 'rgb(213, 150, 221)',
    'z': 'rgb(242, 224, 245)'
}

/**
 * 随机颜色rgb
 *
 * @returns {string}
 */
const randomRGBColor = function () {
    let r = Math.floor(Math.random() * 256);
    let g = Math.floor(Math.random() * 256);
    let b = Math.floor(Math.random() * 256);
    return `rgb(${r},${g},${b})`;
}

/**
 * 随机颜色十六进制值
 *
 * @returns {string}
 */
const randomColor=()=> {
    let str = Math.ceil(Math.random() * 16777215).toString(16);
    if (str.length < 6) {
        str = `0${str}`;
    }
    // return `#${Math.floor(Math.random()*(2<<23)).toString(16)}`;
    return `#${str}`;
}

/**
 * 随机颜色hsl
 *
 * @returns {string}
 */
const randomHSLColor = function () {
    // Hue(色调)。0(或360)表示红色，120表示绿色，
    // 240表示蓝色，也可取其他数值来指定颜色。
    let h = Math.round(Math.random() * 360);
    // Saturation(饱和度)。取值为：0.0% - 100.0%
    let s = Math.round(Math.random() * 100);
    // Lightness(亮度)。取值为：0.0% - 100.0%
    let l = Math.round(Math.random() * 80);
    return `hsl(${h},${s}%,${l}%)`;
}


/**
 * 是否为css合法颜色值
 *
 * @param value
 * @returns {boolean}
 */
const isColor = function (value) {
    let colorReg = /^#([a-fA-F0-9]){3}(([a-fA-F0-9]){3})?$/;
    let rgbaReg = /^[rR][gG][bB][aA]\(\s*((25[0-5]|2[0-4]\d|1?\d{1,2})\s*,\s*){3}\s*(\.|\d+\.)?\d+\s*\)$/;
    let rgbReg = /^[rR][gG][bB]\(\s*((25[0-5]|2[0-4]\d|1?\d{1,2})\s*,\s*){2}(25[0-5]|2[0-4]\d|1?\d{1,2})\s*\)$/;
    let hslReg = /^[hH][sS][lL]\(([0-9]|[1-9][0-9]|[1-3][0-5][0-9]|360)\,(100|[1-9]\d|\d)(.\d{1,2})?%\,(100|[1-9]\d|\d)(.\d{1,2})?%\)$/;

    return colorReg.test(value) || rgbaReg.test(value) || rgbReg.test(value) || hslReg.test(value);
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
    colorArr,
    randomRGBColor,
    randomColor,
    randomHSLColor,
    isColor
}