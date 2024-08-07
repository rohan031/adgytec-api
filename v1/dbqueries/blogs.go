package dbqueries

import "github.com/jackc/pgx/v5"

const CreateBlogItem = `
	INSERT INTO blogs 
	(blog_id, user_id, project_id, title, cover_image, short_text, content, author)
	VALUES 
	(@blogId, @userId, @projectId, @title, @cover, @summary, @content, @author)
`

func CreateBlogItemArgs(
	blogId,
	userId,
	projectId,
	title,
	cover,
	summary,
	content,
	author string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"blogId":    blogId,
		"userId":    userId,
		"projectId": projectId,
		"title":     title,
		"cover":     cover,
		"summary":   summary,
		"content":   content,
		"author":    author,
	}
}

const GetBlogsByProjectId = `
	SELECT blog_id, title, cover_image, short_text, created_at, author
	FROM blogs
	WHERE project_id = @projectId
	ORDER BY created_at DESC
`

func GetBlogsByProjectIdArgs(projectId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"projectId": projectId,
	}
}

const GetBlogById = `
	SELECT blog_id, title, cover_image, short_text, created_at, author, updated_at, content
	FROM blogs
	WHERE blog_id = @blogId
`

func GetBlogsByIdArgs(blogId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"blogId": blogId,
	}
}

const PatchBlogMetadataById = `
	UPDATE blogs 
	SET title=@title, short_text=@summary
	WHERE blog_id=@blogId
`

func PatchBlogMetadataByIdArgs(title, summary, blogId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"title":   title,
		"summary": summary,
		"blogId":  blogId,
	}
}

const DeleteBlogById = `
	DELETE FROM blogs
	WHERE blog_id=@blogId
`

func DeleteBlogByIdArgs(blogId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"blogId": blogId,
	}
}

const PatchBlogCover = `
	WITH cover AS (
		SELECT cover_image as image
		FROM blogs 
		WHERE blog_id = @blogId
	)
	UPDATE blogs
	SET cover_image  = @cover
	WHERE blog_id = @blogId
	RETURNING (
		SELECT image FROM cover
	)
`

func PatchBlogCoverArgs(blogId, cover string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"blogId": blogId,
		"cover":  cover,
	}
}

const PatchBlogContent = `
	UPDATE blogs
	SET content = @content
	WHERE blog_id = @blogId
`

func PatchBlogContentArgs(blogId, content string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"blogId":  blogId,
		"content": content,
	}
}
