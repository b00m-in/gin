package tmpls

import (
        //"html/template"
        "path/filepath"
        "github.com/golang/glog"
        "github.com/gin-gonic/gin"
        "github.com/gin-contrib/multitemplate"
        "b00m.in/tmpl"
        "b00m.in/subs"
)

var (
        Dflt_ctgrs *subs.Category
)

func SetupTemplates(e *gin.Engine, test interface{}) {
        e.Delims("{{", "}}")
        e.SetFuncMap(tmpl.FuncMap)
        //e.LoadHTMLFiles("./templates/base.tmpl", "./templates/body.tmpl")
        r := LoadTemplates(e, "./templates")
        e.HTMLRender = r
        if tp, ok := test.(*subs.Category); ok {
                Dflt_ctgrs = tp
                glog.Infof("%v\n", Dflt_ctgrs)
        }
}

func LoadTemplates(e *gin.Engine, templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/main/*.html")
	if err != nil {
		panic(err.Error())
	}

	sublayouts, err := filepath.Glob(templatesDir + "/layouts/sub/*.html")
	if err != nil {
		panic(err.Error())
	}

	subincludes, err := filepath.Glob(templatesDir + "/includes/sub/*.js")
	if err != nil {
		panic(err.Error())
	}

	subcontents, err := filepath.Glob(templatesDir + "/content/sub/*.html")
	if err != nil {
		panic(err.Error())
	}

	contents, err := filepath.Glob(templatesDir + "/content/main/*.html")
	if err != nil {
		panic(err.Error())
	}
        glog.Infof("layouts: %d, includes: %d\n", len(sublayouts), len(subincludes))

	// Generate our templates map from our layouts/ and content/ directories
	for _, content := range contents {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, content)
                //glog.Infof("files: %v \n", files)
		r.AddFromFilesFuncs(filepath.Base(content), e.FuncMap, files...)
	}
	for _, content := range subcontents {
		layoutCopy := make([]string, len(sublayouts))
		copy(layoutCopy, sublayouts)
                layoutCopy = append(layoutCopy, subincludes...)
		files := append(layoutCopy, content)
		r.AddFromFilesFuncs(filepath.Base(content), e.FuncMap, files...)
	}
	return r

}

