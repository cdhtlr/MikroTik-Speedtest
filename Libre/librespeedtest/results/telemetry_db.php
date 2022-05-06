<?php

define('TELEMETRY_SETTINGS_FILE', 'telemetry_settings.php');

/**
 * @return PDO|false
 */
function getPdo()
{
    if (
        !file_exists(TELEMETRY_SETTINGS_FILE)
        || !is_readable(TELEMETRY_SETTINGS_FILE)
    ) {
        return false;
    }

    require TELEMETRY_SETTINGS_FILE;

    $pdoOptions = [
        PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION
    ];

    try {
		if (!isset($Sqlite_db_file)) {
			return false;
		}

		$pdo = new PDO('sqlite:'.$Sqlite_db_file, null, null, $pdoOptions);

		$pdo->exec('
			CREATE TABLE IF NOT EXISTS `speedtest_users` (
			`id`    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			`timestamp`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			`ip`    text NOT NULL,
			`ua`    text NOT NULL,
			`lang`  text NOT NULL,
			`dl`    text,
			`ul`    text,
			`ping`  text,
			`jitter`        text
			);
		');

		return $pdo;
    } catch (Exception $e) {
        return false;
    }

    return false;
}

/**
 * @return string|false returns the id of the inserted column or false on error
 */
function insertSpeedtestUser($ip, $ua, $lang, $dl, $ul, $ping, $jitter)
{
    $pdo = getPdo();
    if (!($pdo instanceof PDO)) {
        return false;
    }

    try {
        $stmt = $pdo->prepare(
            'INSERT INTO speedtest_users
        (ip,ua,lang,dl,ul,ping,jitter)
        VALUES (?,?,?,?,?,?,?)'
        );
        $stmt->execute([
            $ip, $ua, $lang, $dl, $ul, $ping, $jitter
        ]);
        $id = $pdo->lastInsertId();
    } catch (Exception $e) {
        return false;
    }

    return $id;
}


/**
 * @return array|false
 */
function getLatestSpeedtestUsers()
{
    $pdo = getPdo();
    if (!($pdo instanceof PDO)) {
        return false;
    }

    try {
        $stmt = $pdo->query(
            'SELECT
            timestamp, ip, ua, lang, dl, ul, ping, jitter
            FROM speedtest_users
            ORDER BY timestamp DESC
            LIMIT 8'
        );

        $rows = $stmt->fetchAll(PDO::FETCH_ASSOC);

        foreach ($rows as $i => $row) {
            $rows[$i]['id_formatted'] = $row['id'];
        }
    } catch (Exception $e) {
        return false;
    }

    return $rows;
}
