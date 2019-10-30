package simple

import "html/template"

var serviceTmpl = template.Must(template.New("service").Parse(`
package services

import (
	"{{.PkgName}}/model"
	"github.com/mlogclub/simple"
)

var {{.Name}}Service = &{{.CamelName}}Service {}

type {{.CamelName}}Service struct {
}

func (this *{{.CamelName}}Service) Get(id int64) *model.{{.Name}} {
	ret := &model.{{.Name}}{}
	if err := simple.DB().First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *{{.CamelName}}Service) Take(where ...interface{}) *model.{{.Name}} {
	ret := &model.{{.Name}}{}
	if err := simple.DB().Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *{{.CamelName}}Service) QueryCnd(cnd *simple.SqlCnd) (list []model.{{.Name}}, err error) {
	err = cnd.Exec(simple.DB()).Find(&list).Error
	return
}

func (this *{{.CamelName}}Service) Query(queries *simple.ParamQueries) (list []model.{{.Name}}, paging *simple.Paging) {
	queries.StartQuery(simple.DB()).Find(&list)
	queries.StartCount(simple.DB()).Model(&model.{{.Name}}{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *{{.CamelName}}Service) Create(t *model.{{.Name}}) (*model.{{.Name}}, error) {
	if err := simple.DB().Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (this *{{.CamelName}}Service) Update(t *model.{{.Name}}) error {
	return simple.DB().Save(t).Error
}

func (this *{{.CamelName}}Service) Updates(id int64, columns map[string]interface{}) error {
	return simple.DB().Model(&model.{{.Name}}{}).Where("id = ?", id).Updates(columns).Error
}

func (this *{{.CamelName}}Service) UpdateColumn(id int64, name string, value interface{}) error {
	return simple.DB().Model(&model.{{.Name}}{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (this *{{.CamelName}}Service) Delete(id int64) error {
	return simple.DB().Delete(&model.{{.Name}}{}, "id = ?", id).Error
}

`))

var controllerTmpl = template.Must(template.New("controller").Parse(`
package admin

import (
	"{{.PkgName}}/model"
	"{{.PkgName}}/services"
	"github.com/mlogclub/simple"
	"github.com/kataras/iris"
	"strconv"
)

type {{.Name}}Controller struct {
	Ctx             iris.Context
}

func (this *{{.Name}}Controller) GetBy(id int64) *simple.JsonResult {
	t := services.{{.Name}}Service.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *{{.Name}}Controller) AnyList() *simple.JsonResult {
	list, paging := services.{{.Name}}Service.Query(simple.NewQueryParams(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *{{.Name}}Controller) PostCreate() *simple.JsonResult {
	t := &model.{{.Name}}{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.{{.Name}}Service.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *{{.Name}}Controller) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.{{.Name}}Service.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.{{.Name}}Service.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

`))

var viewIndexTmpl = template.Must(template.New("index.vue").Parse(`
<template>
    <section class="page-container">
        <!--工具条-->
        <el-col :span="24" class="toolbar">
            <el-form :inline="true" :model="filters">
                <el-form-item>
                    <el-input v-model="filters.name" placeholder="名称"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" v-on:click="list">查询</el-button>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="handleAdd">新增</el-button>
                </el-form-item>
            </el-form>
        </el-col>

        <!--列表-->
        <el-table :data="results" highlight-current-row border v-loading="listLoading"
                  style="width: 100%;" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55"></el-table-column>
            <el-table-column prop="id" label="编号"></el-table-column>
            {{range .Fields}}
			<el-table-column prop="{{.CamelName}}" label="{{.CamelName}}"></el-table-column>
            {{end}}
            <el-table-column label="操作" width="150">
                <template slot-scope="scope">
                    <el-button size="small" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
                </template>
            </el-table-column>
        </el-table>

        <!--工具条-->
        <el-col :span="24" class="toolbar">
            <el-pagination layout="total, sizes, prev, pager, next, jumper" :page-sizes="[20, 50, 100, 300]"
                           @current-change="handlePageChange"
                           @size-change="handleLimitChange"
                           :current-page="page.page"
                           :page-size="page.limit"
                           :total="page.total"
                           style="float:right;">
            </el-pagination>
        </el-col>


        <!--新增界面-->
        <el-dialog title="新增" :visible.sync="addFormVisible" :close-on-click-modal="false">
            <el-form :model="addForm" label-width="80px" ref="addForm">
                {{range .Fields}}
				<el-form-item label="{{.CamelName}}">
					<el-input v-model="addForm.{{.CamelName}}"></el-input>
				</el-form-item>
                {{end}}
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click.native="addFormVisible = false">取消</el-button>
                <el-button type="primary" @click.native="addSubmit" :loading="addLoading">提交</el-button>
            </div>
        </el-dialog>

        <!--编辑界面-->
        <el-dialog title="编辑" :visible.sync="editFormVisible" :close-on-click-modal="false">
            <el-form :model="editForm" label-width="80px" ref="editForm">
                <el-input v-model="editForm.id" type="hidden"></el-input>
                {{range .Fields}}
				<el-form-item label="{{.CamelName}}">
					<el-input v-model="editForm.{{.CamelName}}"></el-input>
				</el-form-item>
                {{end}}
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click.native="editFormVisible = false">取消</el-button>
                <el-button type="primary" @click.native="editSubmit" :loading="editLoading">提交</el-button>
            </div>
        </el-dialog>
    </section>
</template>

<script>
  import HttpClient from '../../apis/HttpClient'

  export default {
    name: "List",
    data() {
      return {
        results: [],
        listLoading: false,
        page: {},
        filters: {},
        selectedRows: [],

        addForm: {
          {{range .Fields}}
          '{{.CamelName}}': '',
          {{end}}
        },
        addFormVisible: false,
        addLoading: false,

        editForm: {
          'id': '',
          {{range .Fields}}
          '{{.CamelName}}': '',
          {{end}}
        },
        editFormVisible: false,
        editLoading: false,
      }
    },
    mounted() {
      this.list();
    },
    methods: {
      list() {
        let me = this
        me.listLoading = true
		let params = Object.assign(me.filters, {
          page: me.page.page,
          limit: me.page.limit
        })
        HttpClient.post('/api/admin/{{.KebabName}}/list', params)
          .then(data => {
            me.results = data.results
            me.page = data.page
          })
          .finally(() => {
            me.listLoading = false
          })
      },
      handlePageChange (val) {
        this.page.page = val
        this.list()
      },
      handleLimitChange (val) {
        this.page.limit = val
        this.list()
      },
      handleAdd() {
        this.addForm = {
          name: '',
          description: '',
        }
        this.addFormVisible = true
      },
      addSubmit() {
        let me = this
        HttpClient.post('/api/admin/{{.KebabName}}/create', this.addForm)
          .then(data => {
            me.$message({message: '提交成功', type: 'success'});
            me.addFormVisible = false
            me.list()
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      handleEdit(index, row) {
        let me = this
        HttpClient.get('/api/admin/{{.KebabName}}/' + row.id)
          .then(data => {
            me.editForm = Object.assign({}, data);
            me.editFormVisible = true
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      editSubmit() {
        let me = this
        HttpClient.post('/api/admin/{{.KebabName}}/update', me.editForm)
          .then(data => {
            me.list()
            me.editFormVisible = false
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },

      handleSelectionChange(val) {
        this.selectedRows = val
      },
    }
  }
</script>

<style lang="scss" scoped>

</style>

`))
