/**
 *
 * @Description:
 * @Author: claer
 * @File: time.js
 * @Version: 1.0.0
 * @Time: 2019/9/15 20:11
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */


/**
 * 时间转换工具
 *
 * @param date          时间
 * @param join          年月日之间连接的字符
 * @returns {string}
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