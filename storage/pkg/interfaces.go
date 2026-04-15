package pkg

import (
    "context"
    coreModels "techsupport/core/pkg/models"   
    sysModels  "techsupport/sysinfo/pkg/models" 
    stModels "techsupport/storage/pkg/models"      
)


type Repository interface {

    InsertTicket(ctx context.Context, ticket stModels.TicketRecord) (int64, error)
    InsertTicketDetails(ctx context.Context, ticketID int64, details []coreModels.CalcResult) error
    InsertSystemInfo(ctx context.Context, ticketID int64, sysInfo sysModels.SystemInfo) error

    GetTicketsByClaimant(ctx context.Context, claimantTag string) ([]stModels.TicketRecord, error)
    
    GetTicketDecisionView(ctx context.Context, accTag string) ([]stModels.TicketAgentView, error)

    GetDBRecordByTag(ctx context.Context, accTag string) (coreModels.DBRecord, error)
    GetFullUserProfile(ctx context.Context, accTag string) (map[string]any, error)

    GetSessionHistory(ctx context.Context, accTag string) ([]coreModels.Session, error)
}