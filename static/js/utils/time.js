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
 * @param date 时间
 * @param join 年月日之间连接的字符
 * @returns {string}
 */
const formatTime = (date, join) => {
  let d = [
    date.getFullYear(),
    date.getMonth() + 1,
    date.getDate()
  ];
  let t = [
    date.getHours(),
    date.getMinutes(),
    date.getSeconds()
  ];
  let dateString = d.map((item, i, arr) => {
    // 格式化日期，如月、日、时、分、秒保证为2位数
    item = item.toString();
    return item[1] ? item : '0' + item;
  }).join(join);
  let timeString = t.map((item, i, arr) => {
    // 格式化日期，如月、日、时、分、秒保证为2位数
    return item < 10 ? '0' + item : item;
  }).join(":");
  return dateString + ' ' + timeString;
}

