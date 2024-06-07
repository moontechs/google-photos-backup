DEFINE NAMESPACE IF NOT EXISTS backup;

USE NS backup;
DEFINE DATABASE IF NOT EXISTS backup;

DEFINE TABLE users SCHEMAFULL
	PERMISSIONS
		FOR select, update WHERE id = $auth.id
		FOR create, delete NONE;

DEFINE FIELD email ON users TYPE string ASSERT string::is::email($value);
DEFINE FIELD password ON users TYPE string;

DEFINE INDEX email ON users FIELDS email UNIQUE;

DEFINE SCOPE users SESSION 30d
	SIGNIN (
		SELECT * FROM users WHERE email = $email AND crypto::argon2::compare(password, $password)
	);
	-- SIGNUP (
	-- 	CREATE users CONTENT {
	-- 		email: $email,
	-- 		password: crypto::argon2::generate($password)
	-- 	}
	-- );