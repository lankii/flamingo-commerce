package controller

import (
	"flamingo/core/breadcrumbs"
	"flamingo/core/product/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"log"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type (
	// View demonstrates a product view controller
	View struct {
		responder.ErrorAware    `inject:""`
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		domain.ProductService   `inject:""`

		Template string         `inject:"config:core.product.view.template"`
		Router   *router.Router `inject:""`
	}

	// ProductViewData is used for product rendering
	ProductViewData struct {
		// simple / configurable / configurable_with_variant
		RenderContext       string
		SimpleProduct       domain.SimpleProduct
		ConfigurableProduct domain.ConfigurableProduct
		ActiveVariant       domain.Variant
		VariantSelected     bool
	}
)

// Get Response for Product matching sku param
func (vc *View) Get(c web.Context) web.Response {
	product, err := vc.ProductService.Get(c, c.MustParam1("marketplacecode"))

	// catch error
	if err != nil {
		switch errors.Cause(err).(type) {
		case domain.ProductNotFound:
			return vc.ErrorNotFound(c, err)

		default:
			return vc.Error(c, err)
		}
	}

	var viewData ProductViewData

	// 1. Handle Configurables
	if product.Type() == "configurable" {
		configurableProduct := product.(domain.ConfigurableProduct)
		var activeVariant *domain.Variant

		variantCode, err := c.Param1("variantcode")

		if err == nil {
			log.Println("get variant by " + variantCode)
			activeVariant, _ = configurableProduct.Variant(variantCode)
		}
		if activeVariant == nil {
			// 1.A. No variant selected
			// normalize URL
			urlName := makeUrlTitle(product.BaseData().Title)
			if urlName != c.MustParam1("name") {
				return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "name": urlName})
			}
			viewData = ProductViewData{ConfigurableProduct: configurableProduct, VariantSelected: false, RenderContext: "configurable"}
		} else {
			// 1.B. Variant selected
			// normalize URL
			urlName := makeUrlTitle(activeVariant.BasicProductData.Title)
			if urlName != c.MustParam1("name") {
				return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "variantcode": variantCode, "name": urlName})
			}
			log.Printf("Variant Price %v / %v", activeVariant.ActivePrice.Default, activeVariant.ActivePrice)
			viewData = ProductViewData{ConfigurableProduct: configurableProduct, ActiveVariant: *activeVariant, VariantSelected: true, RenderContext: "configurable_with_activevariant"}
		}

	} else {
		// 2. Handle Simples
		// normalize URL
		urlName := makeUrlTitle(product.BaseData().Title)
		if urlName != c.MustParam1("name") {
			return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "name": urlName})
		}

		simpleProduct := product.(domain.SimpleProduct)
		viewData = ProductViewData{SimpleProduct: simpleProduct, RenderContext: "simple"}
	}

	for _, category := range product.BaseData().CategoryPath {
		breadcrumbs.Add(c, breadcrumbs.Crumb{Title: category})
	}

	breadcrumbs.Add(c, breadcrumbs.Crumb{
		Title: product.BaseData().Title,
		URL:   vc.Router.URL("product.view", router.P{"marketplacecode": product.BaseData().MarketPlaceCode, "name": product.BaseData().Title}).String(),
	})

	return vc.Render(c, vc.Template, viewData)
}

func makeUrlTitle(title string) string {
	newTitle := strings.ToLower(strings.Replace(title, " ", "-", -1))
	newTitle = url.QueryEscape(newTitle)

	return newTitle
}
