/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: log.js
 * @Version: 1.0.0
 * @Time: 2019/9/15 20:26
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */

import time from "./time.js";

// const isDebugEnabled = "production";
const isDebugEnabled = "dev";
const isInfoEnabled = true;
const isErrorEnabled = true;
const isWarnEnabled = true;
const isTraceEnabled = true;

/**
 * 自定义颜色打印日志
 *
 * @param title
 * @param content
 * @param backgroundColor 颜色
 */
const log = (title, content, backgroundColor = "#1475b2") => {
    let i = [
        `%c ${title} %c ${content} `,
        "padding: 1px; border-radius: 3px 0 0 3px; color: #fff; background: ".concat("#606060", ";"),
        `padding: 1px; border-radius: 0 3px 3px 0; color: #fff; background: ${backgroundColor};`
    ];
    return function () {
        let t;
        window.console && "function" === typeof window.console.log && (t = console).log.apply(t, arguments);
    }.apply(void 0, i);
}

log("isDebugEnabled", isDebugEnabled, "#42c02e");
log("isInfoEnabled", isInfoEnabled, "#42c02e");
log("isErrorEnabled", isErrorEnabled, "#42c02e");
log("isWarnEnabled", isWarnEnabled, "#42c02e");
log("isTraceEnabled", isTraceEnabled, "#42c02e");

/**
 * 箭头函数是匿名函数，不能作为构造函数，不能使用new
 *
 * 对日志参数解析
 * 格式为：
 *     logger.info("页面{}，点击第{}行", "App.vue", index);
 *
 * @param log 箭头函数不能绑定arguments，取而代之用rest参数
 * @returns {string}
 */
const getParam = (...log) => {
    if (log.length == 0) {
        return "";
    }
    let params = log[0];
    let parentString = params[0].toString();
    // 正则表达式，如须匹配大小写则去掉i
    let re = eval("/" + "{}" + "/ig");
    // 匹配正则
    let ps = parentString.match(re);

    // 参数个数大于1，并且匹配的个数大于0
    if (params.length > 1 && ps != null) {
        // 移除第一个元素并返回该元素
        params.shift();
        for (let i = 0; i < ps.length; i++) {
            parentString = parentString.replace("{}", params[i]);
        }
        // 把替换后的字符串与参数未替换完的拼接起来
        parentString = parentString + params.slice(ps.length).toString();
        return parentString;
    }
    return params.toString();
}

const debug = (...log) => {
    if (isDebugEnabled) {
        console.log(
            `${time.dateFormat(new Date, "yyyy-MM-dd HH:mm:ss")} %c ${getParam(log)}`,
            'color:red;',
            'font-size:15px;color:red;'
        );
    }
}

const logConcat = (...log) => {
    return `${time.dateFormat(new Date, "yyyy-MM-dd HH:mm:ss")} ${getParam(log)}`;
}

const info = (...log) => {
    if (isInfoEnabled) {
        console.info(logConcat(log));
    }
}

const error = (...log) => {
    if (isErrorEnabled) {
        console.error(logConcat(log));
    }
}
const warn = (...log) => {
    if (isWarnEnabled) {
        console.warn(logConcat(log));
    }
}
const trace = (...log) => {
    if (isTraceEnabled) {
        console.trace(logConcat(log));
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
    debug,
    info,
    error,
    warn,
    trace
}