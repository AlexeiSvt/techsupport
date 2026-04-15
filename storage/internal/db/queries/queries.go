package queries

const (
	InsertSysInfo = `
        INSERT INTO system_info (
            ticket_id, os, platform, arch, kernel, cpu_model, 
            cpu_cores, total_ram, hostname, machine_id, username
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
)

const (
	InsertTicketDetails = `
        INSERT INTO ticket_details (
            ticket_id, name, code, value, weight, result, comment, status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
)

const (
	InsertTicketAndReturnID = `
        INSERT INTO tickets (
            acc_tag, 
            claimant_tag, 
            device_id, 
            final_percentage, 
            knowledge_score, 
            penalty_score, 
            ip_info
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING ticket_id;`
)

const (
    SelectDBRecordByTag = `
        SELECT 
            acc_tag, 
            reg_country, 
            reg_city, 
            first_email, 
            phone, 
            first_device, 
            is_donator, 
            reg_date,
            devices,
            first_transaction,
            user_history
        FROM dbrecord 
        WHERE acc_tag = $1;`
)

const (
	SelectDetailsByTicket = `
        SELECT name, code, result, status 
        FROM ticket_details 
        WHERE ticket_id = $1;`
)

const (
	SelectFullUserProfileByTag = `
        SELECT 
            u.acc_tag, u.first_email, u.is_donator, u.reg_date,
            si.os, si.cpu_model, si.machine_id
        FROM users u
        LEFT JOIN tickets t ON u.acc_tag = t.acc_tag
        LEFT JOIN system_info si ON t.ticket_id = si.ticket_id
        WHERE u.acc_tag = $1
        ORDER BY t.created_at DESC
        LIMIT 1;`
)

const (
	SelectSysInfoByTicket = `
        SELECT os, cpu_model, machine_id 
        FROM system_info 
        WHERE ticket_id = $1;`
)

const (
	SelectTicketsByClaimant = `
        SELECT 
            ticket_id, 
            acc_tag as target_account, 
            final_percentage, 
            created_at 
        FROM tickets 
        WHERE claimant_tag = $1 
        ORDER BY created_at DESC;`
)

const (
    SelectTicketDecisionView = `
        SELECT 
            t.ticket_id, 
            t.created_at as submitted_at, 
            t.updated_at, 
            si.machine_id,
            CASE 
                WHEN CAST(REPLACE(t.final_percentage, '%', '') AS NUMERIC) >= 80 THEN 'ACCEPT'
                WHEN CAST(REPLACE(t.final_percentage, '%', '') AS NUMERIC) >= 50 THEN 'WARNING'
                ELSE 'DENY'
            END as decision
        FROM tickets t
        LEFT JOIN system_info si ON t.ticket_id = si.ticket_id
        WHERE t.acc_tag = $1 
        ORDER BY t.created_at DESC;`
)

const (
    SelectUserSessionHistory = `
        SELECT 
            session_id, 
            session_ip, 
            device_id, 
            asn, 
            country, 
            city, 
            start_time, 
            end_time
        FROM sessions 
        WHERE acc_tag = $1 
        ORDER BY start_time DESC;`
)