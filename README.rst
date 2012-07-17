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
        *gforms.BaseForm
        Title    *gforms.StringField
        Text     *gforms.StringField
        IsPublic *gforms.BoolField
        Image    *gaeforms.BlobField
    }

    func NewArticleForm(article *Article) *ArticleForm {
        title := gforms.NewStringField()
        title.MinLen = 1
        title.MaxLen = 500

        text := gforms.NewTextareaStringField()
        text.MinLen = 1

        isPublic := gforms.NewBoolField()
        isPublic.IsRequired = false
        isPublic.Label = "Is public?"

        image := gaeforms.NewBlobField()

        if article != nil {
            title.SetInitial(article.Title)
            text.SetInitial(article.Text())
            isPublic.SetInitial(article.IsPublic)
        }

        f := &ArticleForm{
            BaseForm: &gforms.BaseForm{},
            Title:    title,
            Text:     text,
            IsPublic: isPublic,
            Image:    image,
        }
        gforms.InitForm(f)

        return f
    }
