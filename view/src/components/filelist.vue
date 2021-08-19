<template>
    <div class="fileview">
        <el-tree class="fileviewLeft"
                 :props="props"
                 @node-click="NodeExpand"
                 :load="loadNode"
                 ref="asyncTree"
                 node-key="path"
                 :lazy="true"
                 accordion
                 highlight-current
                 @node-expand="NodeExpand">
            <span class="custom-tree-node" slot-scope="{ node, data }">
                <span><i :class="'el-icon-folder-opened'" />{{ " " + node.label }}</span>
            </span>
        </el-tree>
        <div class="fileviewright">
            <div class="file-continer-main">
                <div class="file-continer-header">
                    <el-breadcrumb separator-class="el-icon-arrow-right">
                        <template v-for="(v,i) in header">
                            <el-breadcrumb-item><a href="#" @click="HeaderClick(i)">{{ v }}</a></el-breadcrumb-item>
                        </template>
                    </el-breadcrumb>
                </div>
                <div class="file-continer file-list-icon">
                    <div v-for="(v,i) in Data" :key="i" class="file" @dblclick="iconcilck(v.path)" @contextmenu.prevent="openContextMenu($event,v)">
                        <i class="path-ico" :class="v.is_dir ? 'el-icon-folder-opened' : 'el-icon-document'" />
                        <div class="title-type-name"><span class="title">{{v.name}}</span></div>
                    </div>

                </div>
            </div>
        </div>
        <ul :style="{left:left+'px',top:top+'px'}" class="contextmenu" v-show="contextMenuVisible">
            <li @click="closetab(0)">文件删除</li>
            <li @click="closetab(1)">生成副本</li>
            <li @click="closetab(2)">文件下载</li>
        </ul>
    </div>
</template>

<script>
import api from "../request/api.js"

export default {
   data() {
       return {
           props: {label: 'name', children: 'children',path:'path'},
           Data:[],
           header:[],
           contextMenuVisible:false,
           left:0,
           top:0,
           file:{is_dir:true,name:"",path:""},
       };
   },
   methods: {
       loadNode(node, resolve) {
           if (node.level === 0) {
               this.GetDiskList(resolve)
           }else {
               this.GetFileList(node,resolve)
           }

       },
       async GetFileList(node,resolve){
           if (node.data.path === ""){
               resolve([])
               return
           }
           const response = await api.GetFileList(node.data)
           if (response.data && response.data.length > 0){
               node.data.children = response.data
               this.Data = response.data
               this.header = node.data.path.split("\\");

               resolve(response.data.filter(data => data.is_dir ))
           }else {
               this.header = node.data.path.split("\\");
               resolve([])
           }
       },
       async GetDiskList(resolve){
           const response = await api.GetDiskinfo()
           if (response.data.length > 0){
               resolve(response.data)
           }else {
               resolve([])
           }
       },
       NodeExpand(data){
           this.header = data.path.split("\\");
           api.GetFileList(data).then(response =>{
               data.children = response.data
               this.Data = response.data
           })
       },
       iconcilck(data){
           console.log(data)
           let node = this.$refs.asyncTree.getNode(data); // 通过节点id找到对应树节点对象

           node.loaded = false;
           this.Data = []
           node.expand(); // 主动调用展开节点方法，重新查询该节点下的所有子节点
       },
       HeaderClick(e){
           let path = ""
           for (let i= 0; i < e; i++){
               path += this.header[i] + "\\"
           }
           path += this.header[e]

           this.iconcilck(path)
       },
       closetab(type){
           switch (type) {
               case 0:
                   if(this.file.is_dir){
                       this.$message({
                           showClose: true,
                           message: "不支持删除文件夹"
                       })
                   }else {
                       api.DelFile(this.file).then(res => {
                           if (res.data != "OK"){
                               this.$message({
                                   showClose: true,
                                   message: res.data
                               })
                           }else {
                               this.$message({
                                   showClose: true,
                                   message: "文件：" + this.file.name + "删除成功！"
                               })
                               this.HeaderClick(this.header.length-1)
                           }

                       })
                   }
                    break
               case 1:
                   if(this.file.is_dir){
                       this.$message({
                           showClose: true,
                           message: "不支持文件夹"
                       })
                   }else {
                       api.CopyFile(this.file).then(res => {
                           if (res.data != "OK"){
                               this.$message({
                                   showClose: true,
                                   message: res.data
                               })
                           }else {
                               this.$message({
                                   showClose: true,
                                   message: "文件：" + this.file.name + "副本生成完成！"
                               })
                               this.HeaderClick(this.header.length-1)
                           }

                       })
                   }
                   break
               case 2:
                   if(this.file.is_dir){
                       this.$message({
                           showClose: true,
                           message: "不支持文件夹"
                       })
                   }else {
                       //window.open("FileDownload?path="+this.file.path)
                       this.handleDownload("FileDownload?path="+this.file.path, '_self', this.file.name)
                   }
                   break
           }
           console.log(type)
       },
       openContextMenu(e,v){
           console.log(e,v)
           this.file = v
           this.contextMenuVisible = true
           this.left = e.clientX
           this.top = e.clientY + 10
       },
       handleDownload(url, target = '_self', name = 'download') {
           const link = document.createElement('a')
           link.download = name
           link.style.display = 'none'
           link.href = url
           link.target = target
           document.body.appendChild(link)
           link.click()
           document.body.removeChild(link)
       }
   },
    watch:{
        contextMenuVisible() {
            if (this.contextMenuVisible) {
                document.body.addEventListener('click', () => {
                    this.contextMenuVisible = false
                })
            } else {
                document.body.removeEventListener('click', () => {
                    this.contextMenuVisible = false
                })
            }
        },
    },
};
</script>

<style scoped>
.fileview{
        width: -webkit-fill-available;
        height: -webkit-fill-available;
    }

.fileviewLeft{
        overflow:auto;
        width: 300px;
        margin: 20px 0px 0px 20px;
        height: -webkit-fill-available;
        border: 1px solid #fff;
        border-radius:10px;
        padding-top: 20px;
        box-shadow: 0px 0px 5px 5px rgba(0,0,0,0.1);
    }

.fileviewright{
        position: absolute;
        height: -webkit-fill-available;
        left: 320px;
        right: 10px;
        top: 0px;
        margin: 20px 0px 20px 20px;
        border: 1px solid #fff;
        border-radius:10px;
        background: #fff;
        box-shadow: 0px 0px 5px 5px rgba(0,0,0,0.1);
}

.file-continer-main {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    padding-bottom: 50px;
    overflow: auto;
}


.file{
    width: 86px;
    max-height: 145px;
    border: 1px solid transparent;
    padding: 0px;
    box-shadow: 0px 0px 2px rgba(255,255,255,0);
    border-radius: 0;
    filter: none;
    color: #335;
    transition: transform 0.2s;
    text-decoration: none;
    margin-right: 10px;
    margin-bottom: 10px;
    overflow: hidden;
    text-align: center;
    display: inline-block;
    height: auto;
    vertical-align: top;
    position: relative;
}
.file-list-icon .file .title-type-name {
    width: 76px;
    cursor: default;
    text-align: center;
    word-break: break-all;
    font-size: 1.0em;
    margin: 0 auto;
    line-height: 20px;
    padding: 6px 2px 6px 2px;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 3;
}
.file-list-icon .file > .path-ico, .file-list-icon .file > .path-ico .x-item-icon {
    height: 75px;
    width: 75px;
    line-height: 75px;
    font-size: 75px;
}
.file-list-icon .file .title-type-name .title {
    display: block;
}
.file:hover{
    background-color: rgba(0,0,0,0.2);
}

.file-continer-header{
    margin: 10px;
    font-size: 24px;
}

.contextmenu {
    width: 100px;
    margin: 0;
    border: 1px solid #ccc;
    background: #fff;
    z-index: 3000;
    position: fixed;
    list-style-type: none;
    padding: 5px 0;
    border-radius: 4px;
    font-size: 14px;
    color: #333;
    box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.2);
}

.contextmenu li {
    margin: 0;
    padding: 7px 16px;
}

.contextmenu li:hover {
    background: #f2f2f2;
    cursor: pointer;
}
</style>