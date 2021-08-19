<template>
    <div>
       <el-table :data="Data"
                 max-height="600" border style="width: 100%" size="mini">
            <template v-for="(col, index) in Columns">
                <el-table-column
                        v-if="col.slot"
                        :key="index"
                        :prop="col.dataIndex"
                        :label="col.title"
                        :min-width="col.width"
                        :formatter="col.formatter"
                        :sortable ="col.sortable"
                        align="center"
                        show-overflow-tooltip
                >
                    <template v-if="col.search" slot="header">
                        <el-input v-model="search" size="mini" placeholder="输入关键字搜索"></el-input>
                    </template>
                    <template slot-scope="scope">
                        <slot :name="col.slot" :data="scope.row"></slot>
                    </template>
                </el-table-column>
                <el-table-column
                    v-else
                    :key="index" 
                    :prop="col.dataIndex" 
                    :label="col.title"
                    :min-width="col.width" 
                    :formatter="col.formatter"
                    :sortable ="col.sortable"
                    align="center"
                    show-overflow-tooltip
                >
                </el-table-column>
            </template>
        </el-table>
    </div>
</template>

<script>
  export default {
    name:"bartable",
    props:{
        Data:Array,
        Columns:Array,
    },
    data(){
      return{
          search:'',

      }
    },
    methods:{

    }
  }
</script>


<style scoped>

</style>