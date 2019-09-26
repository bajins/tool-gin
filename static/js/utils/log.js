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

// const isDebugEnabled = "production";
const isDebugEnabled = "dev";
const isInfoEnabled = true;
const isErrorEnabled = true;
const isWarnEnabled = true;
const isTraceEnabled = true;


let loggerName = "[" + getCurrAbsPath() + "]";

console.log(
    "%cisDebugEnabled=%c" + `${isDebugEnabled}` +
    ",%cisInfoEnabled=%c" + `${isInfoEnabled}` +
    ",%cisErrorEnabled=%c " + `${isErrorEnabled}` +
    ",%cisWarnEnabled=%c" + `${isWarnEnabled}` +
    ",%cisTraceEnabled=%c" + `${isTraceEnabled}`
    , 'color:#2db7f5;'
    , 'color:red;'
    , 'color:#2db7f5;'
    , 'color:red;'
    , 'color:red;'
    , 'color:red;'
    , 'background:#aaa;color:#bada55;'
    , 'color:red;'
    , 'color:#2db7f5;'
    , 'color:red;'
);

/**
 * 对日志参数解析
 * 格式为：
 *     logger.info("页面{}，点击第{}行", "App.vue", index);
 *
 * @param log
 * @returns {string}
 */
const getParam = (...log) => {
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
        console.log("%c " + loggerName + " %c " + getParam(log), 'color:red;', 'font-size:15px;color:red;');
    }
};

const info = (...log) => {
    if (isInfoEnabled) {
        console.info(loggerName + getParam(log));
    }
};

const error = (...log) => {
    if (isErrorEnabled) {
        console.error(loggerName + getParam(log));
    }
};
const warn = (...log) => {
    if (isWarnEnabled) {
        console.warn(loggerName + getParam(log));
    }
};
const trace = (...log) => {
    if (isTraceEnabled) {
        console.trace(loggerName + getParam(log));
    }
};

/**
 * 将方法、变量暴露出去
 */
export default {
    debug,
    info,
    error,
    warn,
    trace
};