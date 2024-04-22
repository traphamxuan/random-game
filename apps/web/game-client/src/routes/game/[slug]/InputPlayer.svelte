<script lang="ts">
	import { browser } from "$app/environment";
	import { GameStatus, storeGameState } from "$lib/client/stores/game-state";
	import { onMount } from "svelte";
	
  export let serverID: string;

  let timer: NodeJS.Timeout | undefined;
  let userId = ''
  if (browser) {
    userId = localStorage.getItem('userName') || '';
  }
  type ChipAnswer = {
    id: number;
    answer: string;
    state: 'success' | 'error' | 'warning' | 'primary';
  }

  let currentAnswer = '';
  let listChipAnswer: ChipAnswer[] = [];
  let setAnswer = new Set<string>();
  let answerSent = 0;

  const validate = () => {
    currentAnswer = currentAnswer.toLowerCase().replace(/[^a-f0-9]/g, '')
    currentAnswer = currentAnswer.slice(0, 2)
  }

  const processAnswer = async () => {
    if (timer) {
      return;
    }
    while (answerSent < listChipAnswer.length) {
      const chip = listChipAnswer[answerSent];
      await new Promise((ok, failed) => {
        timer = setTimeout(async () => {
          chip.state = 'primary';
          answerSent++;
          listChipAnswer = listChipAnswer;
          const data = await fetch(`/api/game/${serverID}`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              userId,
              gameId: $storeGameState.gameData.id,
              answer: chip.answer
            })
          })

          if (!data.ok) {
            failed(data.statusText)
            return;
          }
          let resp = await data.json();
          const winner = resp['winner'];
          const answer = resp['answer'];
          const rewards = resp['rewards'];
          chip.state = 'error';
          if (winner || answer) {
            const gameData: GameData = resp;
            clearTimeout(timer)
            timer = undefined;
            if (winner === userId) {
              chip.state = 'success';
            }
            storeGameState.stop(gameData);
            listChipAnswer = [];
            setAnswer = new Set();
            answerSent = 0;
            const response = await fetch(`/api/game/${serverID}/state`, {
              method: 'GET',
              headers: {
                'Content-Type': 'application/json'
              }
            })
            if (!response.ok) {
              failed(response.statusText)
              return;
            }
            resp = await response.json();
            storeGameState.start(resp);
          } else {
            storeGameState.run(rewards);
          }
          listChipAnswer = listChipAnswer;
          ok(null);
        }, 500)
      })
    }
    clearTimeout(timer);
    timer = undefined;
    return;
  }

  const onUserEnter = (e: KeyboardEvent) => {
    // Post answer if user type enter
    if (e.key === 'Enter') {
      if (currentAnswer === '') {
        return;
      }
      if (setAnswer.has(currentAnswer)) {
        console.warn('Answer already exist', currentAnswer);
        currentAnswer = '';
        return;
      }
      setAnswer.add(currentAnswer);
      listChipAnswer.push({
        id: listChipAnswer.length + 1,
        answer: currentAnswer,
        state: 'warning'
      });
      listChipAnswer = listChipAnswer;
      currentAnswer = '';
      processAnswer();
    }
  }

  const removeChip = (answer: ChipAnswer) => {
    if (answer.state === 'warning') {
      listChipAnswer = listChipAnswer.filter((chip) => chip.id !== answer.id);
    }
  }

  const autoFill = () => {
    for (let i = 0; i < 255; i++) {
      const answer = i.toString(16).padStart(2, '0');
      if (setAnswer.has(answer)) {
        continue;
      }
      setAnswer.add(answer);
      listChipAnswer.push({
        id: listChipAnswer.length + 1,
        answer,
        state: 'warning'
      });
    }
    listChipAnswer = listChipAnswer;
    processAnswer();
  }

  onMount(() => {
    if ($storeGameState.status !== GameStatus.Playing) {
      fetch(`/api/game/${serverID}/state`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      }).then(async (data) => {
        if (!data.ok) {
          throw new Error(data.statusText)
        }
        const resp = await data.json();
        storeGameState.start(resp);
        return resp;
      })
    }
  });
</script>

<div class="grid grid-cols-5 gap-1">
  <div class="col-start-1 col-span-4">
      <input type="text" class="input w-full" placeholder="Enter your answer..." bind:value={currentAnswer} on:keypress={onUserEnter} on:input={validate}  disabled={$storeGameState.status !== GameStatus.Playing}/>
  </div>
  <div class="col-start0-4 col-span-1">
      <button type="button" class="btn variant-ghost-primary w-full" on:click={() => autoFill()} disabled={$storeGameState.status !== GameStatus.Playing}>Auto Fill</button>
  </div>
</div>
<div>
  {#each Object.values(listChipAnswer) as answer}
    <button class={`chip m-0.5 variant-filled-${answer.state}`} on:click={() => removeChip(answer)} on:keydown={(e) => {if (e.key === 'Enter') removeChip(answer)}}>{answer.answer}</button>
  {/each}
</div>