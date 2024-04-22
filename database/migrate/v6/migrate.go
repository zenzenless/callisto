package v6

// Migrate implements database.Migrator
func (db *Migrator) Migrate() error {
	stmt := `
BEGIN;

DROP TABLE gov_params;
CREATE TABLE gov_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

ALTER TABLE proposal ADD COLUMN metadata TEXT NOT NULL DEFAULT '';
ALTER TABLE proposal DROP COLUMN proposal_route;
ALTER TABLE proposal DROP COLUMN proposal_type;

ALTER TABLE proposal_deposit ADD COLUMN transaction_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE proposal_deposit DROP CONSTRAINT unique_deposit;
ALTER TABLE proposal_deposit ADD CONSTRAINT unique_deposit UNIQUE (proposal_id, depositor_address, transaction_hash);

ALTER TABLE proposal_vote ADD COLUMN weight TEXT NOT NULL DEFAULT '1.0';

ALTER TABLE proposal_vote DROP CONSTRAINT unique_vote;
ALTER TABLE proposal_vote ADD CONSTRAINT unique_vote UNIQUE (proposal_id, voter_address, option);

ALTER TABLE validator_voting_power ALTER COLUMN voting_power TYPE BIGINT USING voting_power::BIGINT;
ALTER TABLE proposal_validator_status_snapshot ALTER COLUMN voting_power TYPE BIGINT USING voting_power::BIGINT;

COMMIT;
`

	_, err := db.SQL.Exec(stmt)

	return err
}
