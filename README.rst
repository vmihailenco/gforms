HTML forms for Golang
=====================

Installation::

    go get github.com/vmihailenco/gforms

Example
=======

Example::

    package blog

    import (
        "github.com/vmihailenco/gforms"
        "github.com/vmihailenco/gforms/gaeforms"
    )

    type ArticleForm struct {
        gforms.BaseForm
        Title    *gforms.StringField `gforms:",req"`
        Text     *gforms.StringField `gforms:",req"`
        IsPublic *gforms.BoolField
        Image    *gaeforms.BlobField
    }

    func NewArticleForm(article *Article) *ArticleForm {
        f := &ArticleForm{}
        gforms.InitForm(f)

        f.Title.MaxLen = 500
        f.IsPublic.Label = "Is public?"

        if article != nil {
            f.Title.SetInitial(article.Title)
            f.Text.SetInitial(article.Text())
            f.IsPublic.SetInitial(article.IsPublic)
        }

        return f
    }
