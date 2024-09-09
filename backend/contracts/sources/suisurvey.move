#[allow(unused_function, unused_use, unused_variable, unused_mut_parameter)]
module admin::suisurvey{
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
    const ESurveyMaxParticipantsExceeded: u64 = 3;
    const EParticipantAlreadyParticipated: u64 = 4;
    const EInsufficientReward: u64 = 5;

    //==============================================================================================
    // Structs
    //==============================================================================================
    public struct State has key {
        id: UID,
        creators: Table<address, Surveys>,
        all_surveys: vector<ID>//<Form_ID>
    }

    public struct Surveys has store{
        creator: address,
        forms: vector<ID> //<Form_ID>
    }

    public struct Form has key{
        id: UID,
        name: String,
        creator: address,
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
        // survey object id
        id: ID,
        blob_id: u256
    }

    public struct SurveyTerminated has copy, drop {
        creator: address,
        // survey name
        name: String,
        // survey object id
        id: ID,
    }

    public struct ParticipationRecorded has copy, drop {
        participant: address,
        // survey name
        name: String,
        // survey object id
        id: ID,
        blob_id: u256
    }

    public struct RewardDistributed has copy, drop {
        creator: address,
        // survey name
        name: String,
        // survey object id
        id: ID,
    }

    //==============================================================================================
    // Init
    //==============================================================================================

    fun init(ctx: &mut TxContext) {
        transfer::share_object(State{id: object::new(ctx), creators: table::new(ctx), all_surveys: vector::empty()});
    }

    //==============================================================================================
    // Entry Functions
    //==============================================================================================

    entry fun create_survey(
        creator: address,
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
        if(!table::contains(&state.creators, creator)){
            table::add(&mut state.creators, creator, Surveys{creator, forms: vector::empty()})
        };
        let surveys = table::borrow_mut(&mut state.creators, creator);
        let uid = object::new(ctx);
        let id = object::uid_to_inner(&uid);
        let form = Form{
            id: uid,
            name,
            creator,
            blob_id,
            expiration: clock::timestamp_ms(clock) + expiration,
            min_participants,
            max_participants,
            participants: vector::empty(),
            responses: table::new(ctx),
            contract_interaction,
            reward,
        };
        transfer::share_object(form);
        vector::push_back(&mut surveys.forms, id);
        vector::push_back(&mut state.all_surveys, id);
        event::emit(SurveyCreated{
            creator,
            name,
            id,
            blob_id
        })
    }

    entry fun terminate_survey(
        form: &mut Form,
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    ){
        let creator = tx_context::sender(ctx);
        assert!(table::contains(&state.creators, creator), ESurveyNotFound);
        let surveys = table::borrow_mut(&mut state.creators, creator);
        form.expiration = clock::timestamp_ms(clock);
        event::emit(SurveyTerminated{
            creator,
            name: form.name,
            id: object::uid_to_inner(&form.id)
        })
    }

    entry fun participate_survey(
        participant: address,
        form: &mut Form,
        creator: address,
        blob_id: u256, //filled data
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    ){
        assert!(table::contains(&state.creators, creator), ESurveyNotFound);
        let surveys = table::borrow_mut(&mut state.creators, creator);
        assert!(form.expiration > clock::timestamp_ms(clock), ESurveyExpired);
        assert!(form.max_participants < vector::length(&form.participants), ESurveyMaxParticipantsExceeded);
        assert!(!vector::contains(&form.participants, &participant), EParticipantAlreadyParticipated);
        let participant_id = vector::length(&form.participants);
        vector::push_back(&mut form.participants, participant);
        table::add(&mut form.responses, participant_id, blob_id);
        event::emit(ParticipationRecorded{
            participant,
            name: form.name,
            id: object::uid_to_inner(&form.id),
            blob_id
        })
    }

    entry fun distribute_reward(
        form: &mut Form,
        name: String,
        total_reward: &mut Coin<SUI>,
        state: &mut State,
        clock: &Clock, //0x6
        ctx: &mut TxContext
    ){
        let creator = tx_context::sender(ctx);
        assert!(table::contains(&state.creators, creator), ESurveyNotFound);
        let surveys = table::borrow_mut(&mut state.creators, creator);
        assert!(clock::timestamp_ms(clock) > form.expiration, ESurveyStillActive);
        assert!(coin::value(total_reward) >= vector::length(&form.participants)*form.reward, EInsufficientReward);
        let mut n = 0;
        while(n < vector::length(&form.participants)){
            transfer::public_transfer(coin::split(total_reward, form.reward, ctx), *vector::borrow(&form.participants, n));
            n = n + 1;
        };
        event::emit(RewardDistributed{
            creator,
            name,
            id: object::uid_to_inner(&form.id),
        })
    }

    //==============================================================================================
    // Getter Functions
    //==============================================================================================

    public fun get_all_surveys_form_ids(
        creator: address,
        state: &State,
        ctx: &mut TxContext
    ): vector<ID>{
        assert!(table::contains(&state.creators, creator), ESurveyNotFound);
        state.all_surveys
    }

    public fun get_all_form_ids_of_single_creator(
        creator: address,
        state: &State,
        ctx: &mut TxContext
    ): vector<ID>{
        assert!(table::contains(&state.creators, creator), ESurveyNotFound);
        let surveys = table::borrow(&state.creators, creator);
        surveys.forms
    }

    public fun get_specific_response(
        participant: address,
        form: &Form,
        ctx: &mut TxContext
    ): u256{ //output: blob_id
        let (_found, index) = vector::index_of(&form.participants, &participant);
        *table::borrow(&form.responses, index)
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