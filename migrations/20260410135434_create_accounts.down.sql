begin;

drop trigger trg_accounts_updated_at on accounts;
drop table if exists accounts;

commit;
