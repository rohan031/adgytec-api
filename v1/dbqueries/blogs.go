package dbqueries

import "github.com/jackc/pgx/v5"

const CreateBlogItem = `
	INSERT INTO blogs 
	(blog_id, user_id, project_id, title, cover_image, short_text, content, author, category_id)
	VALUES 
	(@blogId, @userId, @projectId, @title, @cover, @summary, @content, @author, @categoryId)
`

func CreateBlogItemArgs(
	blogId,
	userId,
	projectId,
	title,
	cover,
	summary,
	content,
	author, categoryId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"blogId":     blogId,
		"userId":     userId,
		"projectId":  projectId,
		"title":      title,
		"cover":      cover,
		"summary":    summary,
		"content":    content,
		"author":     author,
		"categoryId": categoryId,
	}
}

const GetBlogsByProjectId = `
	SELECT b.blog_id, b.title, b.cover_image, b.short_text, b.created_at, b.author, json_build_object('id', c.category_id, 'name', c.category_name) as category
	FROM blogs b
	LEFT JOIN category c
	ON c.category_id = b.category_id
	WHERE b.project_id = @projectId
	ORDER BY b.created_at DESC
`

func GetBlogsByProjectIdArgs(projectId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"projectId": projectId,
	}
}

const GetBlogById = `
	SELECT b.blog_id, b.title, b.cover_image, b.short_text, b.created_at, b.author, b.updated_at, b.content, c.category_name as category
	FROM blogs b
	INNER JOIN category c
	ON c.category_id = b.category_id
	WHERE blog_id = @blogId;
`

func GetBlogsByIdArgs(blogId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"blogId": blogId,
	}
}

const PatchBlogMetadataById = `
	UPDATE blogs 
	SET title=@title, short_text=@summary, category_id=@categoryId
	WHERE blog_id=@blogId
`

func PatchBlogMetadataByIdArgs(title, summary, blogId, categoryId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"title":      title,
		"summary":    summary,
		"blogId":     blogId,
		"categoryId": categoryId,
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
