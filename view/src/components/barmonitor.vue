<template>
<div class="view">
  <div class="viewLeft">
    <div class="leftbox">
      <span class="icon iconfont icon-cpu lefticon"> CPU</span>
      <el-divider />
      <div id="CPUPercent" :style="{width:'300px',height:'200px'}"></div>
    </div>
    <div class="leftbox" >
      <span class="icon iconfont icon-neicun lefticon"> Memory</span>
      <el-divider />
      <el-row style="margin-left:10px">
        <el-col :span="10"><el-progress type="dashboard" :width="100" :percentage="Percentinfo.Memory.usedPercent" /></el-col>
        <el-col :span="14">
          <el-row style="margin: 16px 0px">总大小：{{ parseInt(Percentinfo.Memory.total/1024/1024) }}MB</el-row>
          <el-row>已使用：{{ parseInt(Percentinfo.Memory.used/1024/1024) }}MB</el-row>
        </el-col>
      </el-row>
    </div>
    <div class="leftbox" >
      <span class="icon iconfont icon-wangluoshiyong lefticon"> Network</span>
      <el-divider />
      <div>
        <div id="netBytesSent" :style="{width:'300px',height:'100px'}"></div>
        <div id="netBytesRecv" :style="{width:'300px',height:'100px'}"></div>
      </div>
    </div>
  </div>
  <div class="viewright">
    <el-tabs class="viewtabs" @tab-click="tabChange" v-model="activeName">
      <el-tab-pane name="systeminfo">
        <span slot="label" class="icon iconfont icon-systemInfo"> 系统信息</span>
          <template v-for="(val, key) in SystemInfo" >
            <p v-if="val instanceof Array" :key="key" style="margin-bottom:2px">
              {{key}} :
              <template v-for="sval in val" >
                {{　sval　}}
              </template>
            </p>
            <p v-else :key="key" style="margin-bottom:2px"> {{key}}: {{val}}</p>
          </template>
      </el-tab-pane>
      <el-tab-pane name="process">
        <span slot="label" class="icon iconfont icon-jinchengguanli"> 进程管理</span>
        <bartable :Data="ProcessData" :Columns="ProcessColumns" >
          <span slot="cpu" slot-scope="scope"> {{ scope.data.cpu + "%" }} </span>
          <span slot="memory" slot-scope="scope">
            {{
              (scope.data.memory/1024).toString().replace(/(\d)(?=(?:\d{3})+$)/g, '$1,') +"  KB"
            }}
          </span>
          <span slot="action" slot-scope="scope">
              <el-button size="mini" type="danger" @click="KillProcess(scope.data.pid)">结束进程</el-button>
            </span>
        </bartable>
      </el-tab-pane>
      <el-tab-pane name="file" style="height: 600px">
        <span slot="label" class="icon iconfont icon-bushu"> 文件管理</span>
        <filelist />
      </el-tab-pane>
      <el-tab-pane name="service">
        <span slot="label" class="icon iconfont icon-fuwuguanli"> 服务管理</span>
        <bartable :Data="ServiceData" :Columns="ServiceColumns" >
            <span slot="state" slot-scope="scope">
                <el-tag :type="scope.data.state == 'Running' ? 'danger' : 'info'">{{ scope.data.state == 'Running' ? '正在运行' : '已停止' }}</el-tag>
            </span>
            <span slot="startMode" slot-scope="scope">
              {{ scope.data.startMode == 'Manual' ? '手动' : (scope.data.startMode == 'Disabled' ? '禁用' : '自动') }}
            </span>
        </bartable>
      </el-tab-pane>
      <el-tab-pane name="net">
        <span slot="label" class="icon iconfont icon-wangluozhenduan"> 网络诊断</span>
        <el-row class="ping">
          <div class="cmd" >
            <div id="bash" style="height: 100%;margin: 10px" />
          </div>
        </el-row>
      </el-tab-pane>
    </el-tabs>
  </div>
</div>
</template>

<script>
import api from "../request/api.js"
import bartable from "./bartable.vue"
import filelist from "./filelist.vue"

const ProcessColumns = [
  {title: '进程ID',dataIndex: 'pid',width: "70"},
  {title: '进程名',dataIndex: 'name',ellipsis: true,},
  {title: '进程路径',dataIndex: 'path',ellipsis: true,},
  {title: '用户名',dataIndex: 'username',ellipsis: true,},
  {title: '会话Id',dataIndex: 'sessionid',ellipsis: true,width: "80"},
  {title: '线程数',dataIndex: 'thread_count', sortable:true,width: "80"},
  {title: '句柄数',dataIndex: 'handle_count',sortable:true,width: "80"},
  {title: '内存使用',dataIndex: 'memory',slot:"memory",sortable:true,width: "100"},
  {title: 'CPU使用',dataIndex: 'cpu',slot:"cpu",sortable:true,width: "100"},
  {title: '操作', width: "240",slot:"action",sortable:true},
];
const ServiceColumns = [
  {title: '服务名',dataIndex: 'name',width: "100",},
  {title: '服务显示名',dataIndex: 'caption',width: "100"},
  {title: '描述',dataIndex: 'description',width: "100"},
  {title: '状态',dataIndex: 'state',width: "100",slot:"state"},
  {title: '启动类型',dataIndex: 'startMode',width: "100",slot:"startMode"},
];

export default {
  name: 'barmonitor',
  data(){
    return {
      SystemInfo:null,
      activeName:'systeminfo',
      ProcessColumns,
      ServiceColumns,
      ProcessData:[],
      ServiceData:[],
      myVar:null,
      Percentinfo:{cpu:0,Memory:{total:0,usedPercent:0,used:0}},
      isopen:false,
      pingvalue:1,
      netBytesRecv: {NetChart:null,Bytes:0,netrecv:0},
      netBytesSent:{NetChart:null,Bytes:0,netsend:0},
      netchartCount:0,
      cmdprint:[],
      input:"",
      guid:new Date().getTime(),
      command:"",
      CPUChart:{Chart:null,Count:0},
      socket:null,
    }
  },
  components: {
    bartable,
    filelist
  },
  mounted(){
    this.drawLine()
    this.GetSystemInfo()
    this.GetCpuMemory()
    setInterval( () =>{ this.GetCpuMemory(); }, 1000);
    setInterval( () =>{ this.GetnetBytes(); }, 1000);
  },
  methods: {
    GetSystemInfo (){
      api.GetSystemInfo().then(SystemInfo =>{
        this.SystemInfo = SystemInfo.data
      })
    },
    drawLine(){
      Highcharts.setOptions({
        global: {
            useUTC: false
        }
      })
      const netBytesSentoptions = {
        chart: {
          type: 'area',
          marginRight: 10,
          backgroundColor: '#F8F8F8',
        },
        subtitle :{
          text: '单位：B',
          align: 'right',
        },
        title: null,
        xAxis: {
          visible : false
        },
        yAxis: {
          title:null,
          tickAmount: 2
        },
        tooltip: {
          enabled : false,
        },
        legend: {
          enabled: false
        },
        credits:{
          enabled: false
        },
        exporting:{
          enabled: false
        },
        series: [{
          data: [],
          marker:{//线上数据点
            radius:0,
            lineWidth:0,
            lineColor:'#fba845',
            fillColor:'#fba845',
            states:{
              hover:{
                enabled:false
              }
            }
          }
        }],
      }
      const netBytesRecvoptions = {
        chart: {
          type: 'area',
          marginRight: 10,
          backgroundColor: '#F8F8F8',
        },
        subtitle :{
          text: '单位：B',
          align: 'right',
        },
        title: null,
        xAxis: {
          visible : false
        },
        yAxis: {
          title:null,
          tickAmount: 2
        },
        tooltip: {
          enabled : false,
        },
        legend: {
          enabled: false
        },
        credits:{
          enabled: false
        },
        exporting:{
          enabled: false
        },
        series: [{
          data: [],
          marker:{//线上数据点
            radius:0,
            lineWidth:0,
            lineColor:'#fba845',
            fillColor:'#fba845',
            states:{
              hover:{
                enabled:false
              }
            }
          }
        }]
      }
      const CPUPercentRecvoptions = {
        chart: {
          type: 'area',
          marginRight: 10,
          backgroundColor: '#F8F8F8',
        },
        subtitle :{
          text: 'CPU：%',
          align: 'right',
        },
        title: null,
        xAxis: {
          visible : false
        },
        yAxis: {
          title:null,
          tickAmount: 2
        },
        tooltip: {
          enabled : false,
        },
        legend: {
          enabled: false
        },
        credits:{
          enabled: false
        },
        exporting:{
          enabled: false
        },
        series: [{
          data: [],
          marker:{//线上数据点
            radius:0,
            lineWidth:0,
            lineColor:'#fba845',
            fillColor:'#fba845',
            states:{
              hover:{
                enabled:false
              }
            }
          }
        }]
      }
      this.netBytesSent.NetChart = Highcharts.chart('netBytesSent', netBytesSentoptions);
      this.netBytesRecv.NetChart = Highcharts.chart('netBytesRecv', netBytesRecvoptions);
      this.CPUChart.Chart = Highcharts.chart('CPUPercent', CPUPercentRecvoptions);
    },
    GetProcesslist (){
      api.GetProcesslist().then(Processlist =>{
        this.ProcessData = Processlist.data
      })
    },
    GetServicelist(){
      api.GetServicelist().then(Servicelist =>{
        this.ServiceData = Servicelist.data
      })
    },
    GetCpuMemory(){
      api.GetCpuMemory().then(res =>{
        this.Percentinfo.cpu = Math.floor(res.data.Cpu * 100) / 100
        this.Percentinfo.Memory.total = res.data.Memory.total
        this.Percentinfo.Memory.used = res.data.Memory.used
        this.Percentinfo.Memory.usedPercent = res.data.Memory.usedPercent
        this.CPUChart.Chart.setTitle(null, { text: "CPU：" + this.Percentinfo.cpu + "%"});
        if(this.CPUChart.Count >= 60){
          this.CPUChart.Chart.series[0].addPoint([this.CPUChart.Count, parseInt(this.Percentinfo.cpu)], true, true);
        }else{
          this.CPUChart.Chart.series[0].addPoint([this.CPUChart.Count, parseInt(this.Percentinfo.cpu)], true);
        }
        this.activeLastPointToolip(this.CPUChart.Chart);
        this.CPUChart.Count++
      })
    },
    tabChange(tab){
      switch (tab.name){
        case "systeminfo":
          clearInterval(this.myVar)
          break
        case "process":
          clearInterval(this.myVar)
          this.GetProcesslist()
          this.myVar = setInterval(() =>{ this.GetProcesslist() }, 3000);
          break
        case "file":
          clearInterval(this.myVar)
          break
        case "service":
          clearInterval(this.myVar)
          this.GetServicelist()
          this.myVar = setInterval(() =>{ this.GetServicelist() }, 3000);
          break
        case "net":
          if(this.socket === null){
            this.cmd()
          }
          clearInterval(this.myVar)
          break
      }
    },
    bytestostrings(bytess) {
      var tempint = bytess
      if (tempint < 1024) {
        return tempint + "B"
      } else {
        tempint = tempint / 1024
        if (tempint < 1024) {
          return tempint.toFixed(2) + "KB"
        } else {
          tempint = tempint / 1024
          if (tempint < 1024) {
            return tempint.toFixed(2) + "MB"
          } else {
            tempint = tempint / 1024
            return tempint.toFixed(2) + "GB"
          }
        }
      }
    },
    activeLastPointToolip(chart) {
      var points = chart.series[0].points;
      chart.tooltip.refresh(points[points.length -1]);
    },
    GetnetBytes(){
      api.GetIOCounters().then(res =>{
        if(this.netBytesSent.Bytes != 0 && this.netBytesRecv.Bytes != 0){
          this.netBytesRecv.netrecv = parseInt(res.data.bytesRecv) - this.netBytesRecv.Bytes
          this.netBytesSent.netsend = parseInt(res.data.bytesSent) - this.netBytesSent.Bytes
        }
        this.netBytesRecv.Bytes = parseInt(res.data.bytesRecv)
        this.netBytesSent.Bytes = parseInt(res.data.bytesSent)
        this.netsend(this.netBytesSent)
        this.netrecv(this.netBytesRecv)
        this.netchartCount++
      })
    },
    netrecv (chart) {
      var y = chart.netrecv
      chart.NetChart.setTitle(null, { text: "网络接收：" + this.bytestostrings(chart.netrecv)});
      if(this.netchartCount >= 60){
        chart.NetChart.series[0].addPoint([this.netchartCount, y], true, true);
      }else{
        chart.NetChart.series[0].addPoint([this.netchartCount, y], true);
      }
      this.activeLastPointToolip(chart.NetChart);
    },
    netsend (chart) {
      var y = chart.netsend
      chart.NetChart.setTitle(null, { text: "网络发送：" + this.bytestostrings(chart.netsend)});
      if(this.netchartCount >= 60){
        chart.NetChart.series[0].addPoint([this.netchartCount, y], true, true);
      }else{
        chart.NetChart.series[0].addPoint([this.netchartCount, y], true);
      }
     this.activeLastPointToolip(chart.NetChart);
    },
    PingSet(){
      this.isopen = !this.isopen
      this.cmdprint = []
      switch (this.pingvalue) {
        case 1:
          this.SetExecute("1")
          break
        case 2:
          this.SetExecute("2")
          break
        case 3:
          this.SetExecute("www.baidu.com")
          break
        case 4:
          this.SetExecute(this.input)
          break
      }
    },
    GetExecutePrint(count){
      let m = {Count:count,Guid:this.guid}
      api.GetExecutePrint(m).then(res =>{
        if (res.data == null) {
          clearInterval(this.myVar)
          this.isopen = !this.isopen
          return
        }

        this.cmdprint = res.data
      })
    },
    SetExecute(data){
      let a = 0
      let m = {exec:data,Guid:this.guid}
      api.SetExecute(m).then(res =>{
        if (res.data != "err"){
          this.command = res.data
          this.myVar = setInterval( () =>{ this.GetExecutePrint(a); a++ }, 1000);
        }else {
          this.$message('命令执行失败');
        }

      })
    },
    KillProcess(pid){
      let formData = new FormData();
      formData.append('pid', pid);
      api.KillProcess(formData).then(res => {
        this.$message({
          showClose: true,
          message: res.data
        })
      })
    },
    cmd() {
      let socketUrl = window.location.host + "/pty";
      if (window.location.protocol === "https:") {
        socketUrl = "wss://" + socketUrl;
      } else {
        socketUrl = "ws://" + socketUrl;
      }
      this.socket = new WebSocket(socketUrl);
      const _this = this
      this.socket.onerror = function (e) {
        console.log("socket error", e);
      };

      let timer;
      this.socket.onopen = function (e) {
        window.onresize = function () {
          if(timer != undefined){
            clearTimeout(timer);
          }
          timer = setTimeout(setSize, 500);
        }
        const setSize = function () {
          let initialGeometry = term.proposeGeometry(),
                  cols = initialGeometry.cols,
                  rows = initialGeometry.rows;
          term.resize(cols, rows);
          _this.socket.send(JSON.stringify({ type : "resize", "data" : [cols,rows] }));
        };
        const term = new Terminal({
          cols: 33,
          rows: 60,
          useStyle: true,
          screenKeys: true
        });
        term.open(document.getElementById("bash"));
        term.fit();
        term.on('title', function (title) {
          console.log(title)
        });
        setSize();
        term.on('resize', function (data) {
          _this.socket.send(JSON.stringify({ type : "resize", "data" : [data.cols,data.rows] }));
        });
        term.on('data', function (data) {
          _this.socket.send(JSON.stringify({ type : "data", "data" : data }));
        });
        _this.socket.onmessage = function (msg) {
          term.write(msg.data);
        };
      };
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.view{
  background-color: #F8F8F8;
  width: -webkit-fill-available;
  height: -webkit-fill-available;
}

.viewLeft{
  position: absolute;
  width: 300px;
  height: -webkit-fill-available;
  margin: 0px;
  padding-top: 10px;
}

.viewright{
  position: absolute;
  height: -webkit-fill-available;
  left: 300px;
  right: 10px;
  top: 0px;
  margin: 20px;
  border: 1px solid #fff;
  border-radius:10px;
  background: #fff;
  box-shadow: 0px 0px 5px 5px rgba(0,0,0,0.1); 
}

.viewtabs{
  padding: 10px;
}

.ping {
  padding: 10px;
  height: -webkit-fill-available;
}
.cmd{
  margin:0 auto;
  text-align:left;
  height: 500px;
  background: #0a0a0a;
  border: 1px solid #fff;
  border-radius:10px;
}

.el-divider{
  margin-top: 5px;
  margin-bottom: 12px;
}

.lefticon{
  margin-left:10px;
}
.leftbox{
  margin-top: 30px;
}
</style>
