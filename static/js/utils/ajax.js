/**
 * js封装ajax请求
 * 使用new XMLHttpRequest 创建请求对象,所以不考虑低端IE浏览器(IE6及以下不支持XMLHttpRequest)
 *
 * @Description:
 * @Author: bajins www.bajins.com
 * @File: ajax.js
 * @Version: 1.0.0
 * @Time: 2019/9/12 13:01
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */

/**
 * 注意:请求参数如果包含日期类型.是否能请求成功需要后台接口配合
 * @param settings 请求参数模仿jQuery ajax
 */
const request = (settings = {}) => {
    // 初始化请求参数
    settings = Object.assign({
        url: '', // string
        type: '', // string 'GET' 'POST' 'DELETE'
        method: '',
        dataType: '', // string 期望的返回数据类型:'json' 'text' 'document' ...
        responseType: '',
        async: true, //  boolean true:异步请求 false:同步请求 required
        data: null, // any 请求参数,data需要和请求头Content-Type对应
        contentType: "application/x-www-form-urlencoded; charset=UTF-8",
        headers: {},
        timeout: 1000, // string 超时时间:0表示不设置超时
        beforeSend: (xhr) => {

        },
        success: (result, status, xhr) => {

        },
        error: (xhr, status, error) => {

        },
        complete: (xhr, status) => {

        }
    }, settings);
    settings.type = settings.type || settings.method;
    settings.dataType = settings.dataType || settings.responseType;
    settings.contentType = settings.contentType || settings.headers["Content-Type"] || settings.headers["content-type"];

    // 参数验证
    if (!settings.url) {
        throw new TypeError("ajax请求：url参数不正确");
    }
    if (!settings.type) {
        throw new TypeError("ajax请求：type或method参数不正确");
    }
    if (!settings.dataType) {
        throw new TypeError("ajax请求：dataType或responseType参数不正确");
    }
    if (!settings.method) {
        throw new TypeError("ajax请求：type或method参数不正确");
    }
    // 创建XMLHttpRequest请求对象
    let xhr = new XMLHttpRequest();
    // 请求开始回调函数，对应xhr.loadstart
    xhr.addEventListener('loadstart', e => {
        settings.beforeSend(xhr, e);
    });
    // 请求成功回调函数，对应xhr.onload
    xhr.addEventListener('load', e => {
        const status = xhr.status;
        if ((status >= 200 && status < 300) || status === 304) {
            let result;
            if (xhr.responseType === 'text') {
                result = xhr.responseText;
            } else if (xhr.responseType === 'document') {
                result = xhr.responseXML;
            } else {
                result = xhr.response;
            }
            // 注意:状态码200表示请求发送/接受成功,不表示业务处理成功
            settings.success(result, status, xhr);
        } else {
            settings.error(xhr, status, e);
        }
    });
    // 请求结束，对应xhr.onloadend
    xhr.addEventListener('loadend', e => {
        settings.complete(xhr, xhr.status, e);
    });
    // 请求出错，对应xhr.onerror
    xhr.addEventListener('error', e => {
        settings.error(xhr, xhr.status, e);
    });
    // 请求超时，对应xhr.ontimeout
    xhr.addEventListener('timeout', e => {
        settings.error(xhr, 408, e);
    });

    // 初始化请求
    xhr.open(settings.type, settings.url, settings.async);
    // 设置期望的返回数据类型
    xhr.responseType = settings.dataType || settings.responseType;
    // 设置请求头
    for (const key of Object.keys(settings.headers)) {
        xhr.setRequestHeader(key, settings.headers[key]);
    }
    xhr.setRequestHeader("Content-Type", settings.contentType);

    // 设置超时时间
    if (settings.async && settings.timeout) {
        xhr.timeout = settings.timeout;
    }

    let method = settings.type.toUpperCase();
    // 如果是"简单"请求,则把data参数组装在url上
    if ((method === 'GET' || method === 'DELETE') && settings.data) {
        let paramsStr;
        if ((settings.data).constructor === Object) {
            let paramsArr = [];
            Object.keys(settings.data).forEach(key => {
                paramsArr.push(`${encodeURIComponent(key)}=${encodeURIComponent(settings.data[key])}`);
            });
            paramsStr = paramsArr.join('&');
        } else if ((settings.data).constructor === String) {
            paramsStr = settings.data;
        } else if ((settings.data).constructor === Array) {
            paramsStr = settings.data.join("&");
        } else {
            throw new TypeError("ajax请求：数据类型错误！");
        }
        settings.url += (settings.url.indexOf('?') !== -1) ? paramsStr : '?' + paramsStr;
        xhr.send();
    }
    // 请求参数类型需要和请求头Content-Type对应
    else {
        let ct = settings.contentType.split(";")[0].toLocaleLowerCase();
        if (ct == "application/x-www-form-urlencoded" && (settings.data).constructor === Object) {
            let paramsArr = [];
            Object.keys(settings.data).forEach(key => {
                paramsArr.push(`${encodeURIComponent(key)}=${encodeURIComponent(settings.data[key])}`);
            });
            xhr.send(paramsArr.join('&'));
        } else if (ct == "multipart/form-data" && (settings.data).constructor === Object) {
            let formData = new FormData();
            Object.keys(settings.data).forEach(key => {
                formData.append(key, settings.data[key]);
            });
            xhr.send(formData);
        } else if (ct == "text/plain") {
            if ((settings.data).constructor === String) {
                xhr.send(settings.data);
            } else if ((settings.data).constructor === Array || (settings.data).constructor === Object) {
                xhr.send(JSON.stringify(settings.data));
            } else {
                throw new TypeError("ajax请求：数据类型错误！");
            }
        } else if (ct == "application/json") {
            if ((settings.data).constructor === String) {
                try {
                    JSON.parse(settings.data);
                    xhr.send(settings.data);
                } catch (e) {
                    throw new TypeError("ajax请求：数据类型错误！");
                }
            } else if ((settings.data).constructor === Array || (settings.data).constructor === Object) {
                xhr.send(JSON.stringify(settings.data));
            } else {
                throw new TypeError("ajax请求：数据类型错误！");
            }
        } else {
            throw new TypeError("ajax请求：数据类型错误！");
        }
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
    request
}