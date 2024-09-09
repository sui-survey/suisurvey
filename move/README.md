## SuiSurvey
on-chain survey/polling/voting site. 

### Functions:
1. create_survey(
        name: String,
        expiration: u64, //in ms
        min_participants: u64,
        max_participants: u64,
        blob_id: u256, // empty form data 
        //criterias, if applicable
        contract_interaction: address, //0x0 if N/A
        reward: u64, //per participant, in sui. 0 if N/A
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    )

2. terminate_survey(
        name: String,
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    )

3. participate_survey(
        name: String,
        creator: address,
        blob_id: u256, //filled data 
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    )

4. distribute_reward(
        name: String,
        total_reward: &mut Coin<SUI>,
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    )