import { Link, Outlet } from 'react-router-dom'

export default function Layout() {
  return (
    <div>
      <nav className="header" style={{marginBottom:24}}>
        <div style={{fontSize:18, fontWeight:600}}>OpsAgent</div>
        <div>
          <Link to="/dashboard" style={{color:'#60a5fa'}}>Dashboard</Link>
          <Link to="/alerts" style={{marginLeft:24}}>Alerts</Link>
          <Link to="/assistant" style={{marginLeft:24}}>Assistant</Link>
          <Link to="/audit" style={{marginLeft:24}}>Audit</Link>
          <Link to="/settings" style={{marginLeft:24}}>Settings</Link>
        </div>
      </nav>
      <Outlet />
    </div>
  )
}
