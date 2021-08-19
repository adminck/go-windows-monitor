import axios from 'axios'; // 引入axios

const service = axios.create({
    baseURL: process.env.NODE_ENV === 'development' ? 'http://127.0.0.1:3000' : '',
    timeout: 99999
})

var api = {
    GetProcesslist : (data) => {
        return service({
            url: "/GetProcesslist",
            method: "GET",
            params: data
        })
    },
    GetSystemInfo : (data) => {
        return service({
            url: "/GetSystemInfo",
            method: "GET",
            data: data
        })
    },
    GetDiskinfo : (data) => {
        return service({
            url: "/GetDiskinfo",
            method: "GET",
            data: data
        })
    },
    GetCpuMemory : (data) => {
        return service({
            url: "/GetCpuMemory",
            method: "GET",
            data: data
        })
    },
    GetServicelist : (data) => {
        return service({
            url: "/GetServicelist",
            method: "GET",
            params: data
        })
    },
    GetFileList : (data) => {
        return service({
            url: "/GetFileList",
            method: "GET",
            params: data
        })
    },
    GetIOCounters : (data) => {
        return service({
            url: "/GetIOCounters",
            method: "GET",
            params: data
        })
    },
    SetExecute:(data) => {
        return service({
            url: "/SetExecute",
            method: "GET",
            params: data
        })
    },
    GetExecutePrint:(data) => {
        return service({
            url: "/GetExecutePrint",
            method: "GET",
            params: data
        })
    },
    KillProcess:(data) => {
        return service({
            url: "/KillProcess",
            method: "POST",
            data: data
        })
    },
    DelFile:(data) => {
        return service({
            url: "/DelFile",
            method: "GET",
            params: data
        })
    },
    CopyFile:(data) => {
        return service({
            url: "/CopyFile",
            method: "GET",
            params: data
        })
    },
    FileDownload:(data) => {
        return service({
            url: "/FileDownload",
            method: "GET",
            params: data
        })
    },
}

export default api;