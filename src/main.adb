with Ada.Text_IO;

procedure Main is
begin

   while (True) loop

      Ada.Text_IO.Put("> ");
      declare
         input : String := Ada.Text_IO.Get_Line;
      begin
         Ada.Text_IO.Put_Line (input);
      end;

   end loop;
   null;
end Main;
