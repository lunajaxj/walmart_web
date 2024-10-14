package mode

type ShopMapping struct {
	PID      string `db:"pid"`
	ShopName string `db:"shop_name"`
}
type Category struct {
	CategoryName   string `db:"category_name"`    //类目名称
	CategoryUpName string `db:"category_up_name"` //类目上级
	Children       []*Category
}

type ProductBrands struct {
	Sellers string `db:"sellers"`
	Brands  string `db:"brands"`
	Count   int    `db:"count"`
	Brandss []Brands
}

type ProductBrandss struct {
	Sellers string `db:"sellers"`
	Brandss []int
}

type Brands struct {
	Brands string `db:"brands"`
	Count  int    `db:"count"`
}
type Keyword struct {
	Name string `db:"name"` //关键词名称
	Ids  string `db:"ids"`  //ids

}
type Shopping struct {
	PrId             int    `db:"pr_id"`     //Id值
	Sku              string `db:"sku"`       //sku值
	Price            string `db:"price"`     //当前价格
	XPrice           string `db:"xprice"`    //下架价格
	Inventory        string `db:"inventory"` //库存数量
	IsInventory      bool   //是否库存操作
	IsTakeDown       bool   //是否下架操作
	CenterId         string `db:"center_id"`         //库存CenterID
	Seller           string `db:"seller"`            //卖家
	FloorPrice       string `db:"floor_price"`       //最低价格
	IsActive         bool   `db:"is_active"`         //是否活动
	PromotionsStatus string `db:"promotions_status"` //Active  |  Delete All
	PromoPrice       string `db:"promo_price"`       //活动价格
	PromoStartDate   string `db:"promo_start_date"`  //开始时间
	PromoEndDate     string `db:"promo_end_date"`    //结束时间
	IsUp             bool   `db:"is_up"`             //是否需要更新
	Msg              string `db:"msg"`               //信息
	Status           string
	Status1          string `db:"status1"`          //状态
	Status2          string `db:"status2"`          //状态
	Status3          string `db:"status3"`          //状态
	Status4          string `db:"status4"`          //状态
	Status5          string `db:"status5"`          //状态
	Sales            int    `db:"sales"`            //销量
	Note             string `db:"note"`             //备注
	Name             string `db:"name"`             //名称
	Img              string `db:"img"`              //图片地址
	ShoppingCron     string `db:"shopping_cron"`    //抢购物车cron
	TheShelvesCron   string `db:"the_shelves_cron"` //下架购物车cron
	InventoryCron    string `db:"inventory_cron"`   //加库存cron
	XInventoryCron   string `db:"xinventory_cron"`  //加库存cron
	StatusCron       string //状态cron
	StatusCron1      string `db:"status_cron1"` //状态cron
	StatusCron2      string `db:"status_cron2"` //状态cron
	StatusCron3      string `db:"status_cron3"` //状态cron
	StatusCron4      string `db:"status_cron4"` //状态cron
	StatusCron5      string `db:"status_cron5"` //状态cron
	J9999            string
	J99991           string `db:"j99991"`      //是否最低价包含9999
	J99992           string `db:"j99992"`      //是否最低价包含9999
	J99993           string `db:"j99993"`      //是否最低价包含9999
	J99994           string `db:"j99994"`      //是否最低价包含9999
	J99995           string `db:"j99995"`      //是否最低价包含9999
	UpdateDate       string `db:"update_date"` //最近更新时间
}

type ProductDetails struct {
	Mark         string  `db:"mark"`          //标记
	Remark       string  `db:"remark"`        //备注
	Id           int     `db:"id"`            //id
	Img          string  `db:"img"`           //图片
	CodeType     string  `db:"code_type"`     //商品码类型
	Code         string  `db:"code"`          //商品码值
	Brands       string  `db:"brands"`        //品牌
	Tags         string  `db:"tags"`          //标签
	Title        string  `db:"title"`         //标题
	Rating       string  `db:"rating"`        //评分
	Comments     string  `db:"comments"`      //评论数量
	Price        float64 `db:"price"`         //价格
	Sellers      string  `db:"sellers"`       //卖家
	Distribution string  `db:"distribution"`  //配送
	Variants1    string  `db:"variants1"`     //变体1
	Variants2    string  `db:"variants2"`     //变体2
	VariantsId   string  `db:"variants_id"`   //变体id
	ArrivalTime  string  `db:"arrival_time"`  //到达时间
	Category1    string  `db:"category1"`     //类目1
	Category2    string  `db:"category2"`     //类目2
	Category3    string  `db:"category3"`     //类目3
	Category4    string  `db:"category4"`     //类目4
	Category5    string  `db:"category5"`     //类目5
	Category6    string  `db:"category6"`     //类目6
	Category7    string  `db:"category7"`     //类目7
	CategoryName string  `db:"category_name"` //类目id
	CreateTime   string  `db:"create_time"`   //创建时间
	UpdateTime   string  `db:"update_time"`   //更新时间
	Num          int     `db:"num"`           //count
	StarFrom     string  `db:"star_from"`     //starFrom

}
type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Rule     string `db:"rule"`
}

type Tree struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Last     bool   `json:"last"`
	ParentId string `json:"parentId"`
	Children []Tree `json:"children"`
}

type Mail struct {
	MId        int    `db:"m_id"`
	Seller     string `db:"seller"`
	Count      int    `db:"count"`
	Msg        string `db:"msg"`
	UpdateTime string `db:"update_time"` //更新时间
}

type Log struct {
	Id       int    `db:"id"`       //id
	Classify string `db:"classify"` //类型
	Msg      string `db:"msg"`      //信息
	Val      string `db:"val"`      //值
}

type Activity struct {
	Name string `db:"name"` //活动名
	Date string `db:"date"` //获取时间
	Ids  string `db:"ids"`  //该活动中的商品id
}

type TimedUploads struct {
	TuId       int     `db:"tu_id"`
	Name       string  `db:"name"`
	Genre      string  `db:"genre"`
	File       string  `db:"file"`
	Seller     string  `db:"seller"`    // seller 对应 pid
	ShopName   *string `db:"shop_name"` // 关联查询出来的店铺名称
	Msg        string  `db:"msg"`
	Cron       string  `db:"cron"`
	UpdateTime string  `db:"update_time"`
}

type StockAvailability struct {
	ItemId      string `db:"item_id"`     //id
	SalesUser   string `db:"sales_user"`  //销售账号
	Img         string `db:"img"`         //图片
	CySku       string `db:"cy_sku"`      //仓易sku
	CyName      string `db:"cy_name"`     //仓易产品名
	Gtin        string `db:"gtin"`        //gtin
	PtSku       string `db:"pt_sku"`      //平台sku
	Declaration string `db:"declaration"` //产品英文名
	Num         int64  `db:"num"`         //备货数量
	Warehouse   string `db:"warehouse"`   //发货仓库

	LeadTime   int64   `db:"lead_time"`   //备货天数
	Counts     int64   `db:"counts"`      //15,30,60总数
	TransitNum int64   `db:"transit_num"` //在途数量
	LibraryNum int64   `db:"library_num"` //在库数量
	Weighted   float64 `db:"weighted"`    //加权日均
	Remarks1   string  `db:"remarks1"`    //备注1
	Remarks2   string  `db:"remarks2"`    //备注2

	UpdateDate string `db:"update_date"` //更新时间
}

type StoreInformation struct {
	SID                   int     `db:"sid" comment:"SID"`
	AccountShopName       string  `db:"account_shop_name" comment:"账号店铺名"`
	ShopID                int     `db:"shop_id" comment:"店铺ID"`
	PID                   int     `db:"pid" comment:"PID"`
	OnSaleProductCount    int     `db:"on_sale_product_count" comment:"店铺在售产品数量"`
	HasTargetLink         int     `db:"has_target_link" comment:"有标的链接"`
	BSRLink               int     `db:"bsr_link" comment:"BSR的链接"`
	PPLink                int     `db:"pp_link" comment:"PP的链接"`
	WFSDeliveryCount      int     `db:"wfs_delivery_count" comment:"WFS派送数量"`
	WFSDeliveryPercentage float64 `db:"wfs_delivery_percentage" comment:"WFS配送占比"`
	Price0_10             int     `db:"price_0_10" comment:"价格0-10"`
	Price10_15            int     `db:"price_10_15" comment:"价格10-15"`
	Price15_20            int     `db:"price_15_20" comment:"价格15-20"`
	Price20_40            int     `db:"price_20_40" comment:"价格20-40"`
	Price40_60            int     `db:"price_40_60" comment:"价格40-60"`
	Price60Above          int     `db:"price_60_above" comment:"价格60以上"`
	Reviews0_5            int     `db:"reviews_0_5" comment:"评论0-5"`
	Reviews5_10           int     `db:"reviews_5_10" comment:"评论5-10"`
	Reviews10_20          int     `db:"reviews_10_20" comment:"评论10-20"`
	Reviews20_100         int     `db:"reviews_20_100" comment:"评论20-100"`
	Reviews100Above       int     `db:"reviews_100_above" comment:"评论100以上"`
	Brand                 string  `db:"brand" comment:"品牌"`
	AffiliatedCompany     string  `db:"affiliated_company" comment:"隶属公司"`
	SellerName            string  `db:"seller_name" comment:"公司名"`
	Address               string  `db:"address" comment:"公司地址"`
	Country               string  `db:"country" comment:"国家"`
	Read                  float64 `db:"readx" comment:"评分"`
	SellerReviewsNum      int     `db:"seller_reviews_num" comment:"店铺评论数量"`
	Note1                 string  `db:"note_1" comment:"备注1"`
	Note2                 string  `db:"note_2" comment:"备注2"`
	Note3                 string  `db:"note_3" comment:"备注3"`
	Note4                 string  `db:"note_4" comment:"备注4"`
	Note5                 string  `db:"note_5" comment:"备注5"`
	UpdateDate            string  `db:"update_date" comment:"修改时间"`
}

type ProductSales struct {
	ITEMID          int    `db:"item_id" comment:"item_id"`
	CatalogSellerId string `db:"catalog_seller_id" comment:"店铺id"`
	Img             string `db:"img" comment:"图片"`
	Day01           string `db:"day_01" comment:"最近1天"`
	Day02           string `db:"day_02" comment:"最近2天"`
	Day03           string `db:"day_03" comment:"最近3天"`
	Day04           string `db:"day_04" comment:"最近4天"`
	Day05           string `db:"day_05" comment:"最近5天"`
	Day06           string `db:"day_06" comment:"最近6天"`
	Day07           string `db:"day_07" comment:"最近7天"`
	Day15           string `db:"day_15" comment:"最近15天"`
	Day30           string `db:"day_30" comment:"最近30天"`
	Day60           string `db:"day_60" comment:"最近60天"`
	Day90           string `db:"day_90" comment:"最近90天"`
	CreateDate      string `db:"create_date" comment:"创建时间"`
}
type ProductSalesDay struct {
	ITEMID          int    `db:"item_id" comment:"item_id"`
	CatalogSellerId string `db:"catalog_seller_id" comment:"店铺id"`
	Day0            string `db:"day0" comment:"当天销量"`
	CreateDate      string `db:"create_date" comment:"创建时间,只精确到天,创建数据的时候mysql自动填写"`
}
