CREATE TABLE IF NOT EXISTS `users` (
    `id`         BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '用户ID',
    `uuid`       CHAR(36)     NOT NULL                                COMMENT '全局唯一标识 UUID v4',
    `username`   VARCHAR(32)  NOT NULL                                COMMENT '用户名（3-32位字母数字下划线）',
    `email`      VARCHAR(128) NOT NULL                                COMMENT '邮箱',
    `password`   VARCHAR(255) NOT NULL                                COMMENT 'bcrypt 哈希密码',
    `nickname`   VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '昵称，默认同用户名',
    `avatar`     VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '头像 URL',
    `bio`        VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '个人简介',
    `role`       TINYINT      NOT NULL DEFAULT 1                      COMMENT '角色: 1=普通用户 2=管理员',
    `status`     TINYINT      NOT NULL DEFAULT 1                      COMMENT '状态: 1=正常 2=禁用',
    `last_login` DATETIME                                             COMMENT '最后登录时间',
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                              ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at` DATETIME                                             COMMENT '软删除时间（NULL=未删除）',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_users_uuid`     (`uuid`),
    UNIQUE KEY `uq_users_email`    (`email`),
    UNIQUE KEY `uq_users_username` (`username`),
    KEY        `idx_users_role`    (`role`),
    KEY        `idx_users_status`  (`status`),
    KEY        `idx_users_deleted` (`deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='用户表';


CREATE TABLE IF NOT EXISTS `user_tokens` (
    `id`            BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT   COMMENT '令牌ID',
    `user_id`       BIGINT       UNSIGNED NOT NULL                  COMMENT '关联用户ID',
    `refresh_token` VARCHAR(512) NOT NULL                           COMMENT 'JWT Refresh Token 字符串',
    `user_agent`    VARCHAR(512) NOT NULL DEFAULT ''                COMMENT '客户端 User-Agent',
    `client_ip`     VARCHAR(64)  NOT NULL DEFAULT ''                COMMENT '登录客户端 IP',
    `expires_at`    DATETIME     NOT NULL                           COMMENT 'Token 过期时间',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    PRIMARY KEY (`id`),
    KEY `idx_tokens_user_id`      (`user_id`),
    KEY `idx_tokens_refresh_token`(`refresh_token`(64)),
    KEY `idx_tokens_expires`      (`expires_at`),
    CONSTRAINT `fk_tokens_user`
        FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='用户刷新令牌表（支持多设备登录）';


CREATE TABLE IF NOT EXISTS `categories` (
    `id`            BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '分类ID',
    `name`          VARCHAR(64)  NOT NULL                                COMMENT '分类名称（全局唯一）',
    `slug`          VARCHAR(128) NOT NULL                                COMMENT 'URL slug（全局唯一，用于前台路由）',
    `description`   VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '分类描述',
    `cover`         VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '分类封面图 URL',
    `sort_order`    INT          NOT NULL DEFAULT 0                      COMMENT '排序权重，数值越大越靠前',
    `article_count` INT          NOT NULL DEFAULT 0                      COMMENT '关联文章数量（冗余字段，用于加速列表查询）',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                                 ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`    DATETIME                                             COMMENT '软删除时间（NULL=未删除）',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_categories_name`    (`name`),
    UNIQUE KEY `uq_categories_slug`    (`slug`),
    KEY        `idx_categories_sort`   (`sort_order` DESC),
    KEY        `idx_categories_deleted`(`deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章分类表';

CREATE TABLE IF NOT EXISTS `tags` (
    `id`            BIGINT      UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '标签ID',
    `name`          VARCHAR(64) NOT NULL                                COMMENT '标签名称（全局唯一）',
    `slug`          VARCHAR(128) NOT NULL                               COMMENT 'URL slug（全局唯一）',
    `color`         VARCHAR(16) NOT NULL DEFAULT '#6B7280'              COMMENT '标签展示颜色（Hex 色值，如 #00ADD8）',
    `article_count` INT         NOT NULL DEFAULT 0                      COMMENT '关联文章数量（冗余字段）',
    `created_at`    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at`    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP
                                ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`    DATETIME                                            COMMENT '软删除时间（NULL=未删除）',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_tags_name`   (`name`),
    UNIQUE KEY `uq_tags_slug`   (`slug`),
    KEY        `idx_tags_count` (`article_count` DESC),
    KEY        `idx_tags_deleted`(`deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章标签表';

  CREATE TABLE IF NOT EXISTS `articles` (
    `id`             BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '文章ID',
    `uuid`           CHAR(36)     NOT NULL                                COMMENT '全局唯一标识 UUID v4',
    `author_id`      BIGINT       UNSIGNED NOT NULL                       COMMENT '作者用户ID',
    `category_id`    BIGINT       UNSIGNED                                COMMENT '所属分类ID（NULL=未分类）',
    `title`          VARCHAR(256) NOT NULL                                COMMENT '文章标题',
    `slug`           VARCHAR(300) NOT NULL                                COMMENT 'URL slug（全局唯一，用于前台路由）',
    `summary`        VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '文章摘要（手动填写或自动截取正文前200字）',
    `content`        LONGTEXT     NOT NULL                                COMMENT 'Markdown 原始内容',
    `content_html`   LONGTEXT     NOT NULL                                COMMENT '服务端渲染后的 HTML 内容（可用于 RSS / 缓存）',
    `cover`          VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '封面图 URL',
    `status`         TINYINT      NOT NULL DEFAULT 1                      COMMENT '发布状态: 1=草稿 2=已发布 3=已归档',
    `is_top`         TINYINT      NOT NULL DEFAULT 0                      COMMENT '是否置顶: 0=否 1=是',
    `is_featured`    TINYINT      NOT NULL DEFAULT 0                      COMMENT '是否精选: 0=否 1=是',
    `allow_comment`  TINYINT      NOT NULL DEFAULT 1                      COMMENT '是否允许评论: 0=关闭 1=开启',
    `view_count`     INT          NOT NULL DEFAULT 0                      COMMENT '阅读次数',
    `like_count`     INT          NOT NULL DEFAULT 0                      COMMENT '点赞数（与 article_likes 表冗余）',
    `comment_count`  INT          NOT NULL DEFAULT 0                      COMMENT '评论数（与 comments 表冗余）',
    `bookmark_count` INT          NOT NULL DEFAULT 0                      COMMENT '收藏数（与 article_bookmarks 表冗余）',
    `word_count`     INT          NOT NULL DEFAULT 0                      COMMENT '正文字数（Markdown 去除标记后的字符数）',
    `published_at`   DATETIME                                             COMMENT '发布时间（status=2 时由系统自动设置）',
    `created_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                                  ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`     DATETIME                                             COMMENT '软删除时间（NULL=未删除）',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_articles_uuid`      (`uuid`),
    UNIQUE KEY `uq_articles_slug`      (`slug`),
    KEY        `idx_articles_author`   (`author_id`),
    KEY        `idx_articles_category` (`category_id`),
    KEY        `idx_articles_status`   (`status`),
    -- 前台列表常用查询：已发布 + 按时间降序
    KEY        `idx_articles_pub_time` (`status`, `published_at` DESC),
    -- 置顶排序复合索引
    KEY        `idx_articles_top`      (`is_top` DESC, `published_at` DESC),
    -- 精选文章
    KEY        `idx_articles_featured` (`is_featured`, `status`),
    KEY        `idx_articles_deleted`  (`deleted_at`),

    CONSTRAINT `fk_articles_author`
        FOREIGN KEY (`author_id`)   REFERENCES `users`      (`id`)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT `fk_articles_category`
        FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
        ON DELETE SET NULL
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章主表';


-- 文章-标签 多对多关联表
CREATE TABLE IF NOT EXISTS `article_tags` (
    `article_id` BIGINT   UNSIGNED NOT NULL                  COMMENT '文章ID',
    `tag_id`     BIGINT   UNSIGNED NOT NULL                  COMMENT '标签ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '关联创建时间',

    PRIMARY KEY (`article_id`, `tag_id`),
    -- 反向查询（某标签下所有文章）
    KEY `idx_article_tags_tag` (`tag_id`),

    CONSTRAINT `fk_at_article`
        FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT `fk_at_tag`
        FOREIGN KEY (`tag_id`)     REFERENCES `tags`     (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章标签关联表（多对多）';


CREATE TABLE IF NOT EXISTS `comments` (
    `id`          BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '评论ID',
    `article_id`  BIGINT       UNSIGNED NOT NULL                       COMMENT '所属文章ID',
    `user_id`     BIGINT       UNSIGNED                                COMMENT '评论用户ID（NULL=游客）',
    `parent_id`   BIGINT       UNSIGNED                                COMMENT '父评论ID（NULL=顶级评论）',
    `root_id`     BIGINT       UNSIGNED                                COMMENT '根评论ID，用于聚合同一楼层的所有子评论',
    `reply_to_id` BIGINT       UNSIGNED                                COMMENT '回复的目标评论ID（区分直接回复和楼层内回复）',
    `nickname`    VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '游客昵称（登录用户置空）',
    `email`       VARCHAR(128) NOT NULL DEFAULT ''                     COMMENT '游客邮箱（用于 Gravatar 头像，登录用户置空）',
    `content`     TEXT         NOT NULL                                COMMENT '评论正文（纯文本，前端做 XSS 过滤）',
    `ip`          VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '评论者 IP（用于反垃圾）',
    `user_agent`  VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '客户端 User-Agent',
    `status`      TINYINT      NOT NULL DEFAULT 1                      COMMENT '审核状态: 1=待审核 2=已通过 3=已拒绝 4=垃圾评论',
    `like_count`  INT          NOT NULL DEFAULT 0                      COMMENT '评论点赞数',
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '创建时间',
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                               ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`  DATETIME                                             COMMENT '软删除时间',

    PRIMARY KEY (`id`),
    -- 前台：按文章+审核状态查评论列表
    KEY `idx_comments_article`  (`article_id`, `status`),
    -- 按根评论聚合子评论
    KEY `idx_comments_root`     (`root_id`, `status`),
    -- 按父评论查子评论
    KEY `idx_comments_parent`   (`parent_id`),
    -- 后台：按审核状态筛选
    KEY `idx_comments_status`   (`status`, `created_at` DESC),
    -- 关联用户
    KEY `idx_comments_user`     (`user_id`),
    KEY `idx_comments_deleted`  (`deleted_at`),

    CONSTRAINT `fk_comments_article`
        FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT `fk_comments_user`
        FOREIGN KEY (`user_id`)    REFERENCES `users`    (`id`)
        ON DELETE SET NULL
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章评论表（支持二级嵌套，支持游客评论）';


-- 文章点赞表
CREATE TABLE IF NOT EXISTS `article_likes` (
    `id`         BIGINT   UNSIGNED NOT NULL AUTO_INCREMENT   COMMENT '点赞ID',
    `article_id` BIGINT   UNSIGNED NOT NULL                  COMMENT '文章ID',
    `user_id`    BIGINT   UNSIGNED NOT NULL                  COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',

    PRIMARY KEY (`id`),
    -- 保证同一用户对同一文章只能点赞一次
    UNIQUE KEY `uq_likes_article_user` (`article_id`, `user_id`),
    -- 查询某用户点赞的所有文章
    KEY `idx_likes_user`    (`user_id`),
    KEY `idx_likes_created` (`created_at`),

    CONSTRAINT `fk_likes_article`
        FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT `fk_likes_user`
        FOREIGN KEY (`user_id`)    REFERENCES `users`    (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章点赞表';


-- 文章收藏表
CREATE TABLE IF NOT EXISTS `article_bookmarks` (
    `id`         BIGINT   UNSIGNED NOT NULL AUTO_INCREMENT   COMMENT '收藏ID',
    `article_id` BIGINT   UNSIGNED NOT NULL                  COMMENT '文章ID',
    `user_id`    BIGINT   UNSIGNED NOT NULL                  COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',

    PRIMARY KEY (`id`),
    -- 保证同一用户对同一文章只能收藏一次
    UNIQUE KEY `uq_bookmarks_article_user` (`article_id`, `user_id`),
    -- 查询某用户的收藏列表
    KEY `idx_bookmarks_user`    (`user_id`, `created_at` DESC),
    KEY `idx_bookmarks_article` (`article_id`),

    CONSTRAINT `fk_bookmarks_article`
        FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT `fk_bookmarks_user`
        FOREIGN KEY (`user_id`)    REFERENCES `users`    (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文章收藏表';


  CREATE TABLE IF NOT EXISTS `media` (
    `id`            BIGINT        UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '媒体ID',
    `uploader_id`   BIGINT        UNSIGNED NOT NULL                       COMMENT '上传者用户ID',
    `article_id`    BIGINT        UNSIGNED                                COMMENT '关联文章ID（NULL=独立资源，未绑定文章）',
    `original_name` VARCHAR(256)  NOT NULL DEFAULT ''                     COMMENT '客户端上传时的原始文件名',
    `storage_path`  VARCHAR(512)  NOT NULL                                COMMENT '服务端存储相对路径（如 uploads/2026/03/xxx.jpg）',
    `url`           VARCHAR(512)  NOT NULL                                COMMENT '可访问的公开 URL',
    `mime_type`     VARCHAR(128)  NOT NULL DEFAULT ''                     COMMENT 'MIME 类型，如 image/jpeg image/png image/webp',
    `file_size`     BIGINT        NOT NULL DEFAULT 0                      COMMENT '文件大小（字节）',
    `width`         INT           NOT NULL DEFAULT 0                      COMMENT '图片宽度（px），非图片类型为 0',
    `height`        INT           NOT NULL DEFAULT 0                      COMMENT '图片高度（px），非图片类型为 0',
    `hash`          CHAR(64)      NOT NULL DEFAULT ''                     COMMENT 'SHA256 文件内容哈希（用于去重检测）',
    `created_at`    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '上传时间',
    `updated_at`    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP
                                  ON UPDATE CURRENT_TIMESTAMP             COMMENT '更新时间',
    `deleted_at`    DATETIME                                              COMMENT '软删除时间',

    PRIMARY KEY (`id`),
    -- 按上传者查询资源库
    KEY `idx_media_uploader`  (`uploader_id`, `created_at` DESC),
    -- 按文章查询已绑定的图片
    KEY `idx_media_article`   (`article_id`),
    -- 哈希去重查询
    KEY `idx_media_hash`      (`hash`),
    -- 按文件类型筛选
    KEY `idx_media_mime`      (`mime_type`),
    KEY `idx_media_deleted`   (`deleted_at`),

    CONSTRAINT `fk_media_uploader`
        FOREIGN KEY (`uploader_id`) REFERENCES `users`    (`id`)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT `fk_media_article`
        FOREIGN KEY (`article_id`)  REFERENCES `articles` (`id`)
        ON DELETE SET NULL
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='媒体资源表（图片/文件上传管理）';

  CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id`          BIGINT       UNSIGNED NOT NULL AUTO_INCREMENT        COMMENT '日志ID',
    `user_id`     BIGINT       UNSIGNED NOT NULL                       COMMENT '操作用户ID',
    `module`      VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '所属模块，如 article / category / tag / comment / media / user',
    `action`      VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '操作类型，如 create / update / delete / publish / restore',
    `target_id`   BIGINT       UNSIGNED                                COMMENT '操作目标记录 ID（NULL=无具体目标）',
    `description` VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '操作描述，可包含变更摘要',
    `ip`          VARCHAR(64)  NOT NULL DEFAULT ''                     COMMENT '操作者 IP',
    `user_agent`  VARCHAR(512) NOT NULL DEFAULT ''                     COMMENT '操作者客户端 UA',
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '操作时间',

    PRIMARY KEY (`id`),
    -- 按操作用户查询
    KEY `idx_oplogs_user`   (`user_id`, `created_at` DESC),
    -- 按模块+操作类型筛选
    KEY `idx_oplogs_module` (`module`, `action`),
    -- 按目标记录查询变更历史
    KEY `idx_oplogs_target` (`module`, `target_id`),
    -- 按时间范围查询
    KEY `idx_oplogs_time`   (`created_at` DESC),

    CONSTRAINT `fk_oplogs_user`
        FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='后台管理操作日志表';


  -- 文章归档视图：按年月聚合已发布文章数，驱动前台时间线归档页
CREATE OR REPLACE VIEW `article_archive_view` AS
SELECT
    YEAR(`published_at`)  AS `year`,
    MONTH(`published_at`) AS `month`,
    COUNT(*)              AS `article_count`
FROM `articles`
WHERE `status`     = 2        -- 只统计已发布
  AND `deleted_at` IS NULL    -- 排除软删除
  AND `published_at` IS NOT NULL
GROUP BY
    YEAR(`published_at`),
    MONTH(`published_at`)
ORDER BY
    `year`  DESC,
    `month` DESC;


-- 分类文章统计视图：实时统计各分类下已发布文章数（可作为冗余字段同步依据）
CREATE OR REPLACE VIEW `category_article_count_view` AS
SELECT
    c.`id`                                   AS `category_id`,
    c.`name`                                 AS `category_name`,
    c.`slug`                                 AS `category_slug`,
    COUNT(a.`id`)                            AS `article_count`,
    MAX(a.`published_at`)                    AS `last_published_at`
FROM `categories` c
LEFT JOIN `articles` a
    ON  a.`category_id` = c.`id`
    AND a.`status`       = 2
    AND a.`deleted_at`  IS NULL
WHERE c.`deleted_at` IS NULL
GROUP BY c.`id`, c.`name`, c.`slug`;


-- 标签文章统计视图：实时统计各标签下已发布文章数
CREATE OR REPLACE VIEW `tag_article_count_view` AS
SELECT
    t.`id`                AS `tag_id`,
    t.`name`              AS `tag_name`,
    t.`slug`              AS `tag_slug`,
    t.`color`             AS `tag_color`,
    COUNT(at.`article_id`) AS `article_count`
FROM `tags` t
LEFT JOIN `article_tags` at ON at.`tag_id` = t.`id`
LEFT JOIN `articles`      a
    ON  a.`id`         = at.`article_id`
    AND a.`status`      = 2
    AND a.`deleted_at` IS NULL
WHERE t.`deleted_at` IS NULL
GROUP BY t.`id`, t.`name`, t.`slug`, t.`color`;


-- ── 默认管理员账号 ────────────────────────────────────────────────────────────
-- 初始密码：Admin@2026
-- bcrypt hash (cost=10): $2a$10$GfxMXS0M2Vzrw7I.ehMpHuoufOp1tC1UfnI0QYBYd3B8vx0bsPFQu
-- ⚠️ 请在生产环境首次登录后立即修改密码
INSERT INTO `users`
    (`uuid`, `username`, `email`, `password`, `nickname`, `role`, `status`)
VALUES
    (
        'a0000000-0000-0000-0000-000000000001',
        'admin',
        'admin@antblog.dev',
        '$2a$10$GfxMXS0M2Vzrw7I.ehMpHuoufOp1tC1UfnI0QYBYd3B8vx0bsPFQu',
        'AntBlog Admin',
        2,   -- role=2 管理员
        1    -- status=1 正常
    )
ON DUPLICATE KEY UPDATE `id` = `id`;  -- 幂等：已存在时跳过


-- ── 默认文章分类 ─────────────────────────────────────────────────────────────
INSERT INTO `categories`
    (`name`, `slug`, `description`, `sort_order`)
VALUES
    ('技术',    'tech',          '编程、架构、工具等技术相关文章',     100),
    ('随笔',    'essay',         '生活随笔、读书笔记与个人感悟',        90),
    ('工具',    'tools',         '开发工具、效率软件与资源推荐',         80),
    ('开源',    'open-source',   '开源项目介绍与参与经验分享',           70),
    ('未分类',  'uncategorized', '暂未归入具体分类的文章',                0)
ON DUPLICATE KEY UPDATE `id` = `id`;


-- ── 默认文章标签 ─────────────────────────────────────────────────────────────
INSERT INTO `tags`
    (`name`, `slug`, `color`)
VALUES
    ('Go',          'go',           '#00ADD8'),
    ('Vue',         'vue',          '#42B883'),
    ('TypeScript',  'typescript',   '#3178C6'),
    ('MySQL',       'mysql',        '#4479A1'),
    ('Redis',       'redis',        '#DC382D'),
    ('Docker',      'docker',       '#2496ED'),
    ('Kubernetes',  'kubernetes',   '#326CE5'),
    ('架构设计',    'architecture', '#7C3AED'),
    ('性能优化',    'performance',  '#F59E0B'),
    ('开源',        'open-source',  '#10B981'),
    ('Git',         'git',          '#F05032'),
    ('Linux',       'linux',        '#FCC624')
ON DUPLICATE KEY UPDATE `id` = `id`;