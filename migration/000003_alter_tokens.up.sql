ALTER TABLE tokens DROP CONSTRAINT tokens_pkey;
ALTER TABLE tokens ADD PRIMARY KEY(id,refreshToken);