---


---

<h1 id="廣告投放服務">廣告投放服務</h1>
<h1 id="介紹">介紹</h1>
<p>這是一個簡化的廣告投放服務，提供了兩個 RESTful API，一個用於管理廣告資源，另一個用於根據條件列出符合條件的活躍廣告。該服務使用 Golang 開發，並部署在 Docker 環境中。本文將介紹服務的架構、API 使用方法、參數驗證與錯誤處理、單元測試，以及設計上的一些選擇。</p>
<h2 id="架構概覽">架構概覽</h2>
<ul>
<li>使用 Golang 開發，具有高性能和低內存消耗的特性。</li>
<li>使用 Docker 容器化部署，便於開發、測試和部署。</li>
<li>使用 RESTful API 提供服務。</li>
<li>使用 Redis 進行優化，提高廣告檢索的效率</li>
</ul>
<h2 id="api-使用方法">API 使用方法</h2>
<p>Admin API<br>
<strong>POST /admin/ads/create</strong>：創建廣告資源。<br>
請求範例：<br>
{<br>
“title”: “廣告標題”,<br>
“startAt”: “2024-04-10T00:00:00Z”,<br>
“endAt”: “2024-04-20T00:00:00Z”,<br>
“conditions”: {<br>
“age”: [20, 30, 40],<br>
“gender”: [“M”],<br>
“country”: [“TW”],<br>
“platform”: [“android”, “ios”]<br>
}<br>
}</p>
<ul>
<li>響應：返回創建的廣告資源。</li>
</ul>
<h2 id="投放-api">投放 API</h2>
<p><strong>GET /ads</strong>：列出符合條件的活躍廣告。<br>
-   支持的查詢參數：age、gender、country、platform。<br>
-   支持分頁參數：offset 和 limit。<br>
-   請求範例：<code>/ads?country=TW&amp;limit=10</code><br>
-   響應：返回符合條件的活躍廣告列表。</p>
<p>參數驗證與錯誤處理</p>
<p>對請求的參數進行合理的驗證，包括參數的類型、範圍、必填性等。</p>
<ul>
<li>對可能出現的錯誤進行處理，返回適當的錯誤訊息和 HTTP 狀態碼。</li>
</ul>
<h5 id="單元測試">單元測試</h5>
<ul>
<li>使用 Golang 內建的測試框架進行單元測試。</li>
<li>對每個服務的 API 進行單元測試，確保功能的正確性和穩定性。<br>
-<img src="https://github.com/foxdog1011/Backend_work/blob/master/Untitled1.png?raw=true" alt="enter image description here"></li>
</ul>
<h2 id="設計上的選擇">設計上的選擇</h2>
<ul>
<li>選擇使用 Golang 和 Docker 進行開發和部署，以確保高性能和容易管理。</li>
<li>使用 RESTful API 提供服務，提高了服務的可擴展性和互操作性。</li>
<li>使用參數驗證和錯誤處理，增加了服務的可靠性和安全性。</li>
</ul>
<h2 id="使用方式">使用方式</h2>
<ol>
<li>git clone<br>
<a href="https://github.com/foxdog1011/Backend_work.git">https://github.com/foxdog1011/Backend_work.git</a></li>
<li>部署 Docker 容器:<br>
docker-compose up --build -d</li>
<li>使用 API:<br>
請參考上述的 API 使用方法，使用 POSTman 或任何其他 API 測試工具訪問服務。</li>
</ol>
<h5 id="注意事項">注意事項</h5>
<ul>
<li>請確保 Docker 已正確安裝並運行。</li>
<li>請確保服務正確配置並運行，並滿足性能需求。</li>
<li>請參考存儲庫中的文檔和測試代碼，了解更多關於服務的細節和使用方法。</li>
</ul>
<h5 id="授權">授權</h5>
<p>MIT © 2024 Eason Lin</p>

