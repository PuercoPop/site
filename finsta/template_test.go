package finsta

import (
	"io"
	"testing"

	"github.com/PuercoPop/site/blog"
)

func TestRenderPostList(t *testing.T) {
	type args struct {
		w     io.Writer
		posts []blog.Post
	}
	tt := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Post List",
			args: args{
				w:     nil,
				posts: nil,
			},
			want: "",
		},
		{
			name: "One Post List",
			args: args{
				w: nil,
				posts: []blog.Post{
					{
						Title: "Test",
					},
				},
			},
			want: "",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := &renderer{}
			if err := r.RenderPostList(tc.args.w, tc.args.posts); err != nil {
				t.Errorf("renderer.RenderPostList() error = %v, wantErr %v", err, tc.want)
			}
		})
	}

}
