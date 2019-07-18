package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)


//数据库连接全局变量
//database connection
var db *sql.DB = nil
var dbErr error = nil


//数据库连接信息
//global connection information
const (
	dbhostip="localhost"
	dbusername="root"
	dbpassword="ad123456"
	dbname="piaoyu"
	port="3306"
)


//模板目录
//templates directory
var wd,_ = os.Getwd()


//模板定义
//templates
var (
	meta = wd+"/view/meta.html"
	header = wd+"/view/header.html"
	search = wd+"/view/search.html"
	layout = wd+"/view/tmpl.html"
	indexHome = wd+"/view/index.html"
	slideshow = wd+"/view/slideshow.html"
	tags = wd+"/view/tags.html"
	movielist = wd+"/view/movielist.html"
	coming = wd+"/view/coming.html"
	dailytag = wd+"/view/dailytag.html"
	chosenHome = wd+"/view/chosen.html"
	hotShowHome = wd+"/view/hotshow.html"
	userCenterHome = wd+"/view/usercenter.html"
	rushBuyHome = wd+"/view/rushBuy.html"
	introduceHome = wd+"/view/introduce.html"
	liveShowHome = wd+"/view/liveshow.html"
	mineHome = wd+"/view/mine.html"
	exclusiveHome = wd+"/view/exclusive.html"
	layout2Home = wd+"/view/layout2.html"
	baseline = wd+"/view/baseline.html"
	navibar = wd+"/view/navibar.html"
	last = wd+"/view/last.html"
)


//检查错误
//check errors
func checkErr(err error){
	if err!=nil{
		panic(err)
	}
}

//没有找到相关处理函数时统一由这个处理
//if relative handler not found use this
func notFound(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintf(w,"404 not found")
}


//首页的处理函数
//index handler
func index(w http.ResponseWriter,r *http.Request) {
	//电影资料结构体
	type FilmProfile struct {
		Id int
		FilmId int
		Name string //电影名称
		Image string //封面
		FilmType int //类型
	}

	checkErr(dbErr)
	stmt,err:=db.Prepare("SELECT id,filmId,name,image,filmType FROM filmprofile")
	checkErr(err)
	rows ,err :=stmt.Query()
	checkErr(err)

	data := make([]FilmProfile,0)
	var profile FilmProfile

	var id int
	var filmId int
	var name string
	var image string
	var filmType int

	for rows.Next() {
		err :=rows.Scan(&id,&filmId,&name,&image,&filmType)
		checkErr(err)
		profile.Id = id
		profile.FilmId = filmId
		profile.Name = name
		profile.Image = image
		profile.FilmType = filmType
		data = append(data, profile)
	}

	defer stmt.Close()
	defer rows.Close()
	t,_ := template.ParseFiles(layout,meta,indexHome,header,search,slideshow,tags,movielist,coming,dailytag,baseline,navibar,last)
	t.ExecuteTemplate(w, "layout", data)
}


//精选推荐栏目
//Featured selection
func chosen(w http.ResponseWriter,r *http.Request){
	t,_ := template.ParseFiles(chosenHome,meta,header,search,baseline,navibar,last)
	t.ExecuteTemplate(w, "chosen", "")
}

//全球热映
//hot film
func hotshow(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles(hotShowHome,meta,header,search,baseline,navibar,last)
	t.ExecuteTemplate(w, "hotshow", "")
}

//会员中心
//Membership Center
func usercenter(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles(userCenterHome,meta,header,search,baseline,navibar,last)
	t.ExecuteTemplate(w, "usercenter", "")
}


//超值抢购
//Great value for purchase
func rushBuy(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles(rushBuyHome,meta,header,search,baseline,navibar,last)
	t.ExecuteTemplate(w, "rushBuy", "")
}


//详情介绍
//Detailed introduction
func introduce(w http.ResponseWriter, r *http.Request) {
	filmId := r.FormValue("filmId")

	type Detail struct {
		FilmId int
		Image string
		Name string
		Classify string
		Area string
		Duration int
		ReleaseDate string
		PeopleWTS float64
		Recommend float64
		Introduction string
		Actors string
	}

	type Comment struct {
		NickName string
		Photo string
		GiveScore float64
		CommentContent string
		CommentDate string
		Zan int
	}

	type Data struct {
		FilmDetail Detail
		FilmComment[] Comment
	}

	var detail Detail
	var comment Comment
	var data Data

	commentData := make([]Comment,0)

	// get detailData
	var sql string
	sql = "select b.image,b.name,a.filmId,a.classify,a.area,a.duration,a.releaseDate,a.peopleWTS,a.recommend,a.introduction,c.actors from `filmdetail` a left join `filmprofile` b on a.filmId=b.filmId left join `filmbrief` c on a.filmId=c.filmId where a.filmId=" + filmId
	stmtOne, _ := db.Prepare(sql)
	rowsOne, _ := stmtOne.Query()

	for rowsOne.Next() {
		errOne := rowsOne.Scan(&detail.Image,&detail.Name,&detail.FilmId,&detail.Classify,&detail.Area,&detail.Duration,&detail.ReleaseDate,&detail.PeopleWTS,&detail.Recommend,&detail.Introduction,&detail.Actors)
		checkErr(errOne)
	}

	// get commentData
	sql = "select nickName,photo,giveScore,commentContent,commentDate,zan from `filmComment` where filmId=" + filmId + " order by commentDate desc"
	stmtsTwo, _ := db.Prepare(sql)
	rowsTwo, _ := stmtsTwo.Query()
	for rowsTwo.Next() {
		errsTwo := rowsTwo.Scan(&comment.NickName,&comment.Photo,&comment.GiveScore,&comment.CommentContent,&comment.CommentDate,&comment.Zan)
		checkErr(errsTwo)
		commentData = append(commentData, comment)
	}

	data.FilmDetail = detail
	data.FilmComment = commentData

	defer stmtOne.Close()
	defer rowsOne.Close()
	defer stmtsTwo.Close()
	defer rowsTwo.Close()

	t,_ := template.ParseFiles(introduceHome,meta,header,search,baseline,navibar,last)
	t.ExecuteTemplate(w, "introduce", data)
}


//演出频道
//live show channel
func liveshow(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles(liveShowHome,meta,header,search,baseline,navibar,last)
	t.ExecuteTemplate(w, "liveshow", "")
}


//我的主页
//my home page
func mine(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles(mineHome,meta,header,search,baseline,navibar,last)
	t.ExecuteTemplate(w, "mine", "")
}


//代金券
//coupon
func exclusive(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles(exclusiveHome,meta,header,search,baseline,navibar,last)
	t.ExecuteTemplate(w, "exclusive", "")
}


//搜索、全部列表
//movie list、 search result
func process(w http.ResponseWriter, r *http.Request) {
	type OverView struct {
		Id int
		FilmId int
		Name string
		Image string
		FilmType int
		Score float64
		Director string
		Actors string
	}

	var sql string
	baseSql := "select a.id,a.filmId,a.name,a.image,a.filmType,b.score,b.director,b.actors from filmprofile a left join filmbrief b on a.filmId=b.filmId"
	searchContent := r.FormValue("searchContent")
	fType := r.FormValue("t")
	if searchContent == "" {
		if fType != "1" && fType != "2" {
			sql = baseSql
		} else {
			sql = baseSql + " where a.filmType=" + fType
		}
	} else {
		sql = baseSql + " where a.name like '%" + searchContent + "%'"
	}

	stmt, err := db.Prepare(sql)
	checkErr(err)
	rows ,err :=stmt.Query()
	checkErr(err)

	var id int
	var filmId int
	var name string
	var image string
	var filmType int
	var score float64
	var director string
	var actors string

	data := make([]OverView,0)
	var profile OverView

	for rows.Next() {
		err :=rows.Scan(&id,&filmId,&name,&image,&filmType,&score,&director,&actors)
		checkErr(err)
		profile.Id = id
		profile.FilmId = filmId
		profile.Name = name
		profile.Image = image
		profile.FilmType = filmType
		profile.Score = score
		profile.Director = director
		profile.Actors = actors
		data = append(data, profile)
	}

	defer stmt.Close()
	defer rows.Close()

	t,_ := template.ParseFiles(layout2Home,meta,header,search,baseline,navibar,last)
	t.ExecuteTemplate(w, "layout2", data)
}


//处理ajax请求
//deal ajax request
func dealAjaxRequest(w http.ResponseWriter,r *http.Request){
	type Message struct {
      Names string
      Contents string
      Times int64
    }
    m := Message{"Alice", "Hello world", 1294706395881547001}
    b, _ := json.Marshal(m)
    fmt.Fprintf(w,string(b))
}


//数据库初始化
//database initialization
func InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{dbusername, ":", dbpassword, "@tcp(",dbhostip, ":", port, ")/", dbname, "?charset=utf8"}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	db, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	//验证连接
	if err := db.Ping(); err != nil{
		fmt.Println("opon database fail")
		return
	}
}


//入口函数
//entry function
func main() {

	InitDB()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//注册路由,第一个参数指的是请求路径，第二个参数是一个函数类型，表示这个请求需要处理的事情
	http.HandleFunc("/",notFound)
	http.HandleFunc("/index",index)
	http.HandleFunc("/process",process)
	http.HandleFunc("/introduce",introduce)
	http.HandleFunc("/exclusive",exclusive)
	http.HandleFunc("/rushBuy",rushBuy)
	http.HandleFunc("/chosen",chosen)
	http.HandleFunc("/hotshow",hotshow)
	http.HandleFunc("/usercenter",usercenter)
	http.HandleFunc("/liveshow",liveshow)
	http.HandleFunc("/mine",mine)
	http.HandleFunc("/dealAjaxRequest",dealAjaxRequest)

	error := http.ListenAndServe(":9090",nil)
	if error != nil {
		log.Fatal("ListenAndServe:",error)
	}

	defer db.Close()

}