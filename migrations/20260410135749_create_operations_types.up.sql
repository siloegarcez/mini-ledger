begin;

create table if not exists operations_types (
    operation_id smallint primary key,
    description text not null,
    sign_multiplier smallint not null,
    created_at timestamptz not null default now(),
    constraint chk_operations_types_sign_multiplier check (sign_multiplier in (-1, 1)),
    constraint chk_operations_types_description_no_spaces check (description !~ '\s')
);

revoke update on operations_types from public;

insert into operations_types (operation_id, description, sign_multiplier) values
(1, 'PURCHASE', -1),
(2, 'INSTALLMENT_PURCHASE', -1),
(3, 'WITHDRAWAL', -1),
(4, 'PAYMENT', 1);

commit;
