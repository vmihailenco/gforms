HTML forms for Golang
=====================

Installation::

    go get github.com/vmihailenco/gforms

Example
=======

Example::

    package blog

    import (
        "net/http"

        "github.com/vmihailenco/gforms"
    )

    type ArticleForm struct {
        gforms.BaseForm
        Title    *gforms.StringField `gforms:",req"`
        Text     *gforms.StringField `gforms:",req"`
        IsPublic *gforms.BoolField
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

    func (f *ArticleForm) Populate(article *Article) {
        article.Title = f.Title.Value()
        article.Text = f.Text.Value()
        article.IsPublic = f.IsPublic.Value()
    }

    func CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
        form := NewArticleForm(nil)

        if r.Method == "POST" {
            _ = r.ParseForm()
            if gforms.IsFormValid(form, r.Form) {
                article := &Article{}
                form.PopulateArticle(article)

                if err := SaveArticle(article); err != nil {
                    HandleError(w, err)
                    return
                }

                http.Redirect(w, r, "/articles", http.StatusFound)
                return
            }
        }

        data := struct {
            Form *ArticleForm
        } {
            Form: form,
        }
        RenderTemplate(w, data)
    }

Template::

    <form method="post" class="well article">
      {{render .Form.Title "class" "span6"}}
      {{render .Form.Text "class" "span6"}}
      {{render .Form.IsPublic}}

      <div class="form-actions">
        <button type="submit" class="btn btn-primary">Create New Article</button>
      </div>
    </form>
