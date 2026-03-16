CREATE TABLE users_profile
(
    user_id      UUID PRIMARY KEY,
    username     VARCHAR(50) UNIQUE                     NOT NULL,
    display_name VARCHAR(100),
    avatar_url   VARCHAR(255),
    bio          TEXT,
    is_verified  BOOLEAN                  DEFAULT false NOT NULL,
    is_staff     BOOLEAN                  DEFAULT false NOT NULL,
    is_private   BOOLEAN                  DEFAULT false NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE follows
(
    follower_id  UUID                                   NOT NULL REFERENCES users_profile (user_id) ON DELETE CASCADE,
    following_id UUID                                   NOT NULL REFERENCES users_profile (user_id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    PRIMARY KEY (follower_id, following_id)
);

CREATE TABLE post_types
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO post_types (name)
VALUES ('news'),
       ('user_post'),
       ('investment_idea');

CREATE TABLE posts
(
    id                   BIGSERIAL PRIMARY KEY,
    author_id            UUID                                         NOT NULL REFERENCES users_profile (user_id) ON DELETE CASCADE,
    type_id              INT                                          NOT NULL REFERENCES post_types (id),
    content              TEXT                                         NOT NULL,
    media_attachments    JSONB,
    related_asset_ticker VARCHAR(10),
    status               VARCHAR(20)              DEFAULT 'published' NOT NULL,
    likes_count          INT                      DEFAULT 0           NOT NULL,
    comments_count       INT                      DEFAULT 0           NOT NULL,
    is_pinned            BOOLEAN                  DEFAULT false       NOT NULL,
    is_edited            BOOLEAN                  DEFAULT false       NOT NULL,
    created_at           TIMESTAMP WITH TIME ZONE DEFAULT NOW()       NOT NULL,
    updated_at           TIMESTAMP WITH TIME ZONE DEFAULT NOW()       NOT NULL,
    deleted_at           TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_posts_author_id ON posts (author_id);
CREATE INDEX idx_posts_asset_ticker ON posts (related_asset_ticker);

CREATE TABLE post_likes
(
    post_id    BIGINT                                 NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    user_id    UUID                                   NOT NULL REFERENCES users_profile (user_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    PRIMARY KEY (post_id, user_id)
);

CREATE TABLE comments
(
    id                BIGSERIAL PRIMARY KEY,
    post_id           BIGINT                                       NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    user_id           UUID                                         NOT NULL REFERENCES users_profile (user_id) ON DELETE CASCADE,
    parent_comment_id BIGINT REFERENCES comments (id) ON DELETE CASCADE,
    content           TEXT                                         NOT NULL,
    status            VARCHAR(20)              DEFAULT 'published' NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT NOW()       NOT NULL,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT NOW()       NOT NULL,
    deleted_at        TIMESTAMP WITH TIME ZONE
);

CREATE TABLE reports
(
    id          BIGSERIAL PRIMARY KEY,
    reporter_id UUID                                    NOT NULL REFERENCES users_profile (user_id) ON DELETE CASCADE,
    target_type VARCHAR(20)                             NOT NULL,
    target_id   BIGINT                                  NOT NULL,
    reason      TEXT                                    NOT NULL,
    status      VARCHAR(20)              DEFAULT 'open' NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()  NOT NULL
);

CREATE TABLE chat_types
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO chat_types (name)
VALUES ('private'),
       ('group'),
       ('support_channel');

CREATE TABLE chats
(
    id              UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    type_id         INT                                    NOT NULL REFERENCES chat_types (id),
    title           VARCHAR(100),
    avatar_url      VARCHAR(255),
    last_message_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE chat_members
(
    chat_id              UUID                                      NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    user_id              UUID                                      NOT NULL REFERENCES users_profile (user_id) ON DELETE CASCADE,
    role                 VARCHAR(20)              DEFAULT 'member' NOT NULL,
    joined_at            TIMESTAMP WITH TIME ZONE DEFAULT NOW()    NOT NULL,
    last_read_message_id BIGINT,
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages
(
    id                      BIGSERIAL PRIMARY KEY,
    chat_id                 UUID                                   NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    sender_id               UUID                                   NOT NULL REFERENCES users_profile (user_id) ON DELETE CASCADE,
    reply_to_message_id     BIGINT                                 REFERENCES messages (id) ON DELETE SET NULL,
    content                 TEXT,
    media_attachments       JSONB,

    is_transfer             BOOLEAN                  DEFAULT false NOT NULL,
    transfer_amount         NUMERIC(20, 0),
    transfer_currency       VARCHAR(3),
    idempotency_key         VARCHAR(255) UNIQUE,
    transfer_transaction_id UUID,
    transfer_status         VARCHAR(20),

    is_edited               BOOLEAN                  DEFAULT false NOT NULL,
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at              TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_messages_chat_id ON messages (chat_id, id DESC);