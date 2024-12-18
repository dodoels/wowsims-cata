package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/common/cata"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

var ItemSetGladiatorsBattlegear = core.NewItemSet(core.ItemSet{
	ID:   909,
	Name: "Gladiator's Battlegear",
	Bonuses: map[int32]core.ApplySetItemEffect{
		2: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()
			bonusValue := 70.0

			character.MakeCallbackEffectForSetBonus(setName, 2, core.CustomSetBonusCallbackConfig{
				OnGain: func(sim *core.Simulation, _ *core.Aura) {
					// If Sim is undefined ItemSwap is disabled so we can add this statically
					if sim == nil {
						character.AddStat(stats.Strength, bonusValue)
					} else {
						character.AddStatDynamic(sim, stats.Strength, bonusValue)
					}
				},
				OnExpire: func(sim *core.Simulation, _ *core.Aura) {
					character.AddStatDynamic(sim, stats.Strength, -bonusValue)
				},
			})
		},
		4: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()
			bonusValue := 90.0

			character.MakeCallbackEffectForSetBonus(setName, 2, core.CustomSetBonusCallbackConfig{
				OnGain: func(sim *core.Simulation, _ *core.Aura) {
					// If Sim is undefined ItemSwap is disabled so we can add this statically
					if sim == nil {
						character.AddStat(stats.Strength, bonusValue)
					} else {
						character.AddStatDynamic(sim, stats.Strength, bonusValue)
					}
				},
				OnExpire: func(sim *core.Simulation, _ *core.Aura) {
					character.AddStatDynamic(sim, stats.Strength, -bonusValue)
				},
			})
		},
	},
})

var ItemSetEarthenWarplate = core.NewItemSet(core.ItemSet{
	Name: "Earthen Warplate",
	Bonuses: map[int32]core.ApplySetItemEffect{
		2: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()

			setBonusDep := character.AddDynamicMod(core.SpellModConfig{
				ClassMask:  SpellMaskBloodthirst | SpellMaskMortalStrike,
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.05,
			})

			character.MakeCallbackEffectForSetBonus(setName, 2, core.CustomSetBonusCallbackConfig{
				OnGain: func(sim *core.Simulation, _ *core.Aura) {
					setBonusDep.Activate()
				},
				OnExpire: func(sim *core.Simulation, _ *core.Aura) {
					setBonusDep.Deactivate()
				},
			})
		},
		4: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()
			actionID := core.ActionID{SpellID: 90294}

			apDep := make([]*stats.StatDependency, 3)
			for i := 1; i <= 3; i++ {
				apDep[i-1] = character.NewDynamicMultiplyStat(stats.AttackPower, 1.0+float64(i)*0.01)
			}

			buff := character.RegisterAura(core.Aura{
				Label:     "Rage of the Ages",
				ActionID:  actionID,
				Duration:  30 * time.Second,
				MaxStacks: 3,
				OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
					// Example from DK Death Eater
					if oldStacks > 0 {
						character.DisableDynamicStatDep(sim, apDep[oldStacks-1])
					}
					if newStacks > 0 {
						character.EnableDynamicStatDep(sim, apDep[newStacks-1])
					}
				},
			})

			character.MakeProcTriggerAuraForSetBonus(setName, 4, core.ProcTrigger{
				Name:           "Rage of the Ages Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskOverpower | SpellMaskRagingBlow,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					buff.Activate(sim)
					buff.AddStack(sim)
				},
			}, nil)
		},
	},
})

var ItemSetEarthenBattleplate = core.NewItemSet(core.ItemSet{
	Name: "Earthen Battleplate",
	Bonuses: map[int32]core.ApplySetItemEffect{
		2: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()

			setBonusDep := character.AddDynamicMod(core.SpellModConfig{
				ClassMask:  SpellMaskShieldSlam,
				Kind:       core.SpellMod_DamageDone_Flat,
				FloatValue: 0.05,
			})

			character.MakeCallbackEffectForSetBonus(setName, 2, core.CustomSetBonusCallbackConfig{
				OnGain: func(sim *core.Simulation, _ *core.Aura) {
					setBonusDep.Activate()
				},
				OnExpire: func(sim *core.Simulation, _ *core.Aura) {
					setBonusDep.Deactivate()
				},
			})
		},
		4: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()

			setBonusDep := character.AddDynamicMod(core.SpellModConfig{
				ClassMask:  SpellMaskShieldWall,
				Kind:       core.SpellMod_Cooldown_Multiplier,
				FloatValue: -0.5,
			})

			character.MakeCallbackEffectForSetBonus(setName, 2, core.CustomSetBonusCallbackConfig{
				OnGain: func(sim *core.Simulation, _ *core.Aura) {
					setBonusDep.Activate()
				},
				OnExpire: func(sim *core.Simulation, _ *core.Aura) {
					setBonusDep.Deactivate()
				},
			})
		},
	},
})

var ItemSetMoltenGiantWarplate = core.NewItemSet(core.ItemSet{
	Name: "Molten Giant Warplate",
	Bonuses: map[int32]core.ApplySetItemEffect{
		2: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()
			actionID := core.ActionID{SpellID: 99233}

			talentReduction := time.Duration(character.Talents.BoomingVoice*3) * time.Second

			buff := character.RegisterAura(core.Aura{
				Label:    "Burning Rage",
				ActionID: actionID,
				Duration: 12*time.Second - talentReduction,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= 1.1
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] /= 1.1
				},
			})

			character.MakeProcTriggerAuraForSetBonus(setName, 2, core.ProcTrigger{
				Name:           "Burning Rage Trigger",
				ActionID:       actionID,
				ClassSpellMask: SpellMaskShouts,
				Callback:       core.CallbackOnCastComplete,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					buff.Activate(sim)
				},
			}, nil)
		},
		4: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()

			actionID := core.ActionID{SpellID: 99237}

			fieryAttackActionID := core.ActionID{} // actual ID = 99237
			switch character.Spec {
			case proto.Spec_SpecArmsWarrior:
				fieryAttackActionID.SpellID = 12294
			case proto.Spec_SpecFuryWarrior:
				fieryAttackActionID.SpellID = 85288
			}

			fieryAttack := character.RegisterSpell(core.SpellConfig{
				ActionID:    fieryAttackActionID.WithTag(3),
				SpellSchool: core.SpellSchoolFire,
				ProcMask:    core.ProcMaskEmpty,
				Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,

				CritMultiplier:   character.DefaultMeleeCritMultiplier(),
				DamageMultiplier: 1,
				ThreatMultiplier: 1,

				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
					spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				},
			})

			character.MakeProcTriggerAuraForSetBonus(setName, 4, core.ProcTrigger{
				Name:           "Fiery Attack Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskMortalStrike | SpellMaskRagingBlow,
				ProcChance:     0.3,
				Outcome:        core.OutcomeLanded,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					fieryAttack.Cast(sim, result.Target)
				},
			}, nil)
		},
	},
})

var ItemSetMoltenGiantBattleplate = core.NewItemSet(core.ItemSet{
	Name: "Molten Giant Battleplate",
	Bonuses: map[int32]core.ApplySetItemEffect{
		2: func(agent core.Agent, _ string) {
			character := agent.(WarriorAgent).GetWarrior()

			cata.RegisterIgniteEffect(&character.Unit, cata.IgniteConfig{
				ActionID:           core.ActionID{SpellID: 23922}.WithTag(3), // actual 99240
				DisableCastMetrics: true,
				DotAuraLabel:       "Combust",
				IncludeAuraDelay:   true,

				ProcTrigger: core.ProcTrigger{
					Name:           "Combust",
					Callback:       core.CallbackOnSpellHitDealt,
					ClassSpellMask: SpellMaskShieldSlam,
					Outcome:        core.OutcomeLanded,
				},

				DamageCalculator: func(result *core.SpellResult) float64 {
					return result.Damage * 0.2
				},
			})
		},
		4: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()

			aura := character.RegisterAura(core.Aura{
				Label:    "T12 4P Bonus",
				ActionID: core.ActionID{SpellID: 99242},
				Duration: 10 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					character.PseudoStats.BaseParryChance += 0.06
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					character.PseudoStats.BaseParryChance -= 0.06
				},
			})

			character.MakeCallbackEffectForSetBonus(setName, 4, core.CustomSetBonusCallbackConfig{
				OnGain: func(sim *core.Simulation, _ *core.Aura) {
					if sim != nil {
						aura.Activate(sim)
					}
				},
				OnExpire: func(sim *core.Simulation, _ *core.Aura) {
					aura.Deactivate(sim)
				},
			})

		},
	},
})

var ItemSetColossalDragonplateBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Colossal Dragonplate Battlegear",
	Bonuses: map[int32]core.ApplySetItemEffect{
		2: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()

			mod := character.AddDynamicMod(core.SpellModConfig{
				ClassMask:  SpellMaskHeroicStrike,
				Kind:       core.SpellMod_PowerCost_Flat,
				FloatValue: -10,
			})

			actionID := core.ActionID{SpellID: 105860}
			buffAura := character.RegisterAura(core.Aura{
				Label:    "Volatile Outrage",
				ActionID: actionID,
				Duration: 15 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					mod.Activate()
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					mod.Deactivate()
				},
			})

			character.MakeProcTriggerAuraForSetBonus(setName, 2, core.ProcTrigger{
				Name:           "Volatile Outrage Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnCastComplete,
				ClassSpellMask: SpellMaskInnerRage,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					buffAura.Activate(sim)
				},
			}, nil)
		},
		4: func(agent core.Agent, setName string) {
			warrior := agent.(WarriorAgent).GetWarrior()

			actionID := core.ActionID{SpellID: 108126}
			procCS := warrior.RegisterSpell(core.SpellConfig{
				ActionID:    actionID,
				SpellSchool: core.SpellSchoolPhysical,
				Flags:       core.SpellFlagNoOnDamageDealt | core.SpellFlagNoOnCastComplete,
				Cast: core.CastConfig{
					DefaultCast: core.Cast{
						GCD: 0,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					warrior.ColossusSmashAuras.Get(target).Activate(sim)
				},
			})

			// TODO (4.3): Check if this cares that the hit landed
			warrior.MakeProcTriggerAuraForSetBonus(setName, 4, core.ProcTrigger{
				Name:           "Warrior T13 4P Bloodthirst Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskBloodthirst,
				ProcChance:     0.06,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procCS.Cast(sim, result.Target)
				},
			}, nil)

			warrior.MakeProcTriggerAuraForSetBonus(setName, 4, core.ProcTrigger{
				Name:           "Warrior T13 4P Mortal Strike Trigger",
				ActionID:       actionID,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskMortalStrike,
				ProcChance:     0.13,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					procCS.Cast(sim, result.Target)
				},
			}, nil)
		},
	},
})

var ItemSetColossalDragonplateArmor = core.NewItemSet(core.ItemSet{
	Name: "Colossal Dragonplate Armor",
	Bonuses: map[int32]core.ApplySetItemEffect{
		2: func(agent core.Agent, setName string) {
			character := agent.(WarriorAgent).GetWarrior()
			actionID := core.ActionID{SpellID: 105909}
			duration := time.Second * 6

			shieldAmt := 0.0
			shieldAura := character.NewDamageAbsorptionAura("Shield of Fury"+character.Label, actionID, duration, func(unit *core.Unit) float64 {
				return shieldAmt
			})

			character.MakeProcTriggerAuraForSetBonus(setName, 2, core.ProcTrigger{
				Name:           "Shield of Fury Trigger" + character.Label,
				Callback:       core.CallbackOnSpellHitDealt,
				ClassSpellMask: SpellMaskRevenge,
				Outcome:        core.OutcomeLanded,
				ProcChance:     1,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Target != character.CurrentTarget {
						return
					}

					shieldAmt = result.Damage * 0.2
					if shieldAmt > 1 {
						shieldAura.Activate(sim)
					}
				},
			}, &core.CustomSetBonusCallbackConfig{
				OnExpire: func(sim *core.Simulation, aura *core.Aura) {
					shieldAmt = 0
					shieldAura.Deactivate(sim)
					aura.Deactivate(sim)
				},
			})
		},
		4: func(agent core.Agent, _ string) {
			// TODO: Implement this, turns Shield Wall into a raid buff
		},
	},
})

func init() {
}
