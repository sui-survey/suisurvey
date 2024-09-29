#[allow(unused_function, unused_use, unused_variable, unused_mut_parameter)]
module admin::voteverse{
    use sui::event;
    use sui::clock::{Self, Clock};
    use std::string::{Self, String};
    use sui::table::{Self, Table};
    use sui::coin::{Self, Coin};
    use sui::sui::SUI;

    #[test_only]
    use sui::test_scenario::{Self, ctx};
    #[test_only]
    use sui::test_utils::assert_eq;

    //==============================================================================================
    // Constants
    //==============================================================================================
    
    
    //==============================================================================================
    // Error codes
    //==============================================================================================
    /// Survey does not exist
    const ESurveyNotFound: u64 = 0;
    const ESurveyExpired: u64 = 1;
    const ESurveyStillActive: u64 = 2;
    const ESurveyNameTaken: u64 = 3;
    const ESurveyMaxParticipantsExceeded: u64 = 4;
    const EInsufficientReward: u64 = 5;

    //==============================================================================================
    // Structs 
    //==============================================================================================
    public struct State has key {
        id: UID,
        creators: Table<address, Surveys>
    }

    public struct Surveys has store{
        owner: address,
        forms: Table<String, Form> //<form_name, Form>
    }

    public struct Form has store{
        name: String,
        blob_id: u256,
        expiration: u64, //in ms
        min_participants: u64,
        max_participants: u64,
        participants: vector<address>,
        responses: Table<u64, u256>, // <participant_id (vector index), blob_id>
        //criterias, if applicable
        contract_interaction: address, //0x0 if N/A
        reward: u64 //per participant, in sui. 0 if N/A
    }

    //==============================================================================================
    // Event Structs 
    //==============================================================================================

    public struct SurveyCreated has copy, drop {
        creator: address,
        // survey name
        name: String,
        blob_id: u256
    }

    public struct SurveyTerminated has copy, drop {
        creator: address,
        // survey name
        name: String,
    }

    public struct ParticipationRecorded has copy, drop {
        participant: address,
        // survey name
        name: String,
        blob_id: u256
    }

    public struct RewardDistributed has copy, drop {
        creator: address,
        // survey name
        name: String,
    }

    //==============================================================================================
    // Init
    //==============================================================================================

    fun init(ctx: &mut TxContext) {
        transfer::share_object(State{id: object::new(ctx), creators: table::new(ctx)});
    }

    //==============================================================================================
    // Entry Functions 
    //==============================================================================================

    entry fun create_survey(
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
    ){
        let creator = tx_context::sender(ctx);
        if(!table::contains(&state.creators, creator)){
            table::add(&mut state.creators, creator, Surveys{owner: creator, forms: table::new(ctx)})
        };
        let surveys = table::borrow_mut(&mut state.creators, creator);
        assert!(!table::contains(&surveys.forms, name), ESurveyNameTaken);
        let form = Form{
            name,
            blob_id,
            expiration: clock::timestamp_ms(clock) + expiration,
            min_participants,
            max_participants,
            participants: vector::empty(),
            responses: table::new(ctx),
            contract_interaction,
            reward,
        };
        table::add(&mut surveys.forms, name, form);
        event::emit(SurveyCreated{
            creator,
            name,
            blob_id
        })
    }

    entry fun terminate_survey(
        name: String,
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    ){
        let creator = tx_context::sender(ctx);
        assert!(table::contains(&state.creators, creator), ESurveyNotFound);
        let surveys = table::borrow_mut(&mut state.creators, creator);
        assert!(table::contains(&surveys.forms, name), ESurveyNotFound);
        let form = table::borrow_mut(&mut surveys.forms, name);
        form.expiration = clock::timestamp_ms(clock);
        event::emit(SurveyTerminated{
            creator,
            name
        })
    }

    entry fun participate_survey(
        name: String,
        creator: address,
        blob_id: u256, //filled data 
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    ){
        let participant = tx_context::sender(ctx);
        assert!(table::contains(&state.creators, creator), ESurveyNotFound);
        let surveys = table::borrow_mut(&mut state.creators, creator);
        assert!(table::contains(&surveys.forms, name), ESurveyNotFound);
        let form = table::borrow_mut(&mut surveys.forms, name);
        assert!(form.expiration > clock::timestamp_ms(clock), ESurveyExpired);
        assert!(form.max_participants < vector::length(&form.participants), ESurveyMaxParticipantsExceeded);
        let participant_id = vector::length(&form.participants);
        vector::push_back(&mut form.participants, participant);
        table::add(&mut form.responses, participant_id, blob_id);
        event::emit(ParticipationRecorded{
            participant,
            name,
            blob_id
        })
    }

    entry fun distribute_reward(
        name: String,
        total_reward: &mut Coin<SUI>,
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    ){
        let creator = tx_context::sender(ctx);
        assert!(table::contains(&state.creators, creator), ESurveyNotFound);
        let surveys = table::borrow_mut(&mut state.creators, creator);
        assert!(table::contains(&surveys.forms, name), ESurveyNotFound);
        let form = table::borrow_mut(&mut surveys.forms, name);
        assert!(clock::timestamp_ms(clock) > form.expiration, ESurveyStillActive);
        assert!(coin::value(total_reward) >= vector::length(&form.participants)*form.reward, EInsufficientReward);
        let mut n = 0;
        while(n < vector::length(&form.participants)){
            transfer::public_transfer(coin::split(total_reward, form.reward, ctx), *vector::borrow(&form.participants, n));
            n = n + 1;
        };
        event::emit(RewardDistributed{
            creator,
            name
        })
    }

    //==============================================================================================
    // Getter Functions 
    //==============================================================================================

    public fun get_num_participants(
        name: String,
        creator: address,
        state: &mut State,
        ctx: &mut TxContext
    ): u64{
        assert!(table::contains(&state.creators, creator), ESurveyNotFound);
        let surveys = table::borrow_mut(&mut state.creators, creator);
        assert!(table::contains(&surveys.forms, name), ESurveyNotFound);
        let form = table::borrow_mut(&mut surveys.forms, name);
        vector::length(&form.participants)
    }

    //==============================================================================================
    // Helper Functions 
    //==============================================================================================

    fun num_to_string(num: u64): String {
        let mut num_vec = vector::empty<u8>();
        let mut n = num;
        if (n == 0) {
            vector::push_back(&mut num_vec, 48);
        } else {
            while (n != 0) {
                let mod = n % 10 + 48;
                vector::push_back(&mut num_vec, (mod as u8));
                n = n / 10;
            };
        };

        vector::reverse(&mut num_vec);
        string::utf8(num_vec)
    }

    //==============================================================================================
    // Tests 
    //==============================================================================================
    #[test]
    fun test_init_success() {
        let module_owner = @0xa;

        let mut scenario_val = test_scenario::begin(module_owner);
        let scenario = &mut scenario_val;

        {
            init(test_scenario::ctx(scenario));
        };
        let tx = test_scenario::next_tx(scenario, module_owner);
        let expected_events_emitted = 0;
        let expected_created_objects = 1;
        assert_eq(
            test_scenario::num_user_events(&tx), 
            expected_events_emitted
        );
        assert_eq(
            vector::length(&test_scenario::created(&tx)),
            expected_created_objects
        );
        test_scenario::end(scenario_val);
    }
}